### 测试运行文档

---

#### Attention!!!!!
	
在智能实现中，涉及到一下三类分隔符：

	var itemSp = "\n"
	var listSp = "!@#$"
	var minSp = "^&*"
	
在（）中以，分割的数据，在存储过程中都已itemSp进行分割；
在list中的不同项以listSp进行分割；list中单个子项，如果包含多个字段，以minSp分割。

例如（A，B，[(a,b),(c,d)]）的存储格式为:
	
	A + itemSp + B + itemSp + a + minSp + b + listSp + c + minSp + d


---

#### dev环境配置

	CORE_CHAINCODE_ID_NAME=mycc CORE_PEER_ADDRESS=0.0.0.0:7051 ./mychaincode

#### chaincode部署

如果打开了安全模式，需要登录
	
	peer network login jim -p 6avZQLwcUe9b
	
然后进行部署
	
	peer chaincode deploy -u jim -n mycc -c '{"Args": ["init"]}'
	
#### chaincode调用 （CLI）

##### Invoke

###### register

传入参数为userId，完成注册。

注册用户xiaoming

	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "register", "Args": ["xiaoming"]}'
	
注册医院his1和his2
	
	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "register", "Args": ["his1"]}'
	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "register", "Args": ["his2"]}'

###### add

添加数据，需要传入的数据项有以下几项：
	
	（recordId，ownerId，providerId，accessInfo，permission[dataItem,  query, hash, deadline]）

其中，permission中需要存入ownerId的所有权，包含以下字段：
	
	dataItem,  query, hash, deadline
	dataItem：对应数据的表头信息
	query：查询语句
	hash：数据的hash值
	deadline：数据的失效时间
	
例如添加一条xiaoming的数据，存储于his1，获取方式为”ip:port“，该条记录中包含以下字段：name、sex、more，查询语句为”select name,sex,more from xx;"，得到结果hash只为“123123”，数据的时限为：2020/10/10

	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "add", "Args": ["00002","xiaoming","his1","ip:port","name-sex-more","select name,sex,more from xx;","123123","2020/10/10"]}'
	
该操作之后，经添加一条recordId数据在数据，可以通过getRecord方法查询。

###### updateAccessInfo

修改recordId的AccessInfo，该操作由医院提起，需要传入参数为（recordId，providerId, accessInfo）

例如his1需要对recordId为00001的数据进行修改，将accessInfo更改为“http://ip:port”

	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "updateAccessInfo", "Args": ["00002","his1","http://ip:port"]}'

	
###### updatePermission

传入参数(recordId，ownerId，[consumer, dataItem,  query, hash, deadline]）

（该操作有recordId的owner发起）

例如xiaoming要将recordId为00002的数据共享给his2，共享的数据项为name和sex，其查询语句为“select name,sex from xx;”，hash值为“110”，时效为2020/10/10

	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "updatePermission", "Args": ["00002","xiaoming","his2","name-sex","select name,sex from xx;","110","2020/10/10"]}'


###### confirmsummary

需要传入数据为（ownerId，recordId，op）

	ownerId为操作者id，recordId为需要进行操作recordId
	其中op为操作类型：接收为“ac”，拒绝为“re”

用于确认提交到自己summary list中的record。例如his1添加一条数据xiaoming的数据后，其summary状态被改变，通过getSummary查询到改变的recordId后，对其进行接收或者拒绝操作。

xiaoming对recordId为00001的数据进行确认操作：

	peer chaincode invoke -u jim -l golang -n mycc -c '{"Function": "confirmSummary", "Args": ["xiaoming", "00001", "ac"]}'
	
	

##### Query

###### getSummary

用于获取用户的summary数据，传入参数为userId。

	peer chaincode query -u jim -l golang -n mycc -c '{"Function": "getSummary", "Args": ["xiaoming"]}'

返回数据按分隔符分开，依次为：flag、count、recordIdlist。

###### getRecord

用于获取单条recordId记录的数据，传入参数为userId和recordId

	peer chaincode query -u jim -l golang -n mycc -c '{"Function": "getRecord", "Args": ["xiaoming","00002"]}'
	
###### havePermission

传入参数为（userId, recordId， query）

	peer chaincode query -u jim -l golang -n mycc -c '{"Function": "havePermission", "Args": ["his2","00002","select name,sex from xx;"]}'

	

