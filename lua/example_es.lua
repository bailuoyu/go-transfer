local ops = require("esOps") --加载elasticsearch操作模块

local row = ops.rawRow()  --当前数据库的一行数据,table类型，key为列名称
local action = ops.rawAction()  --当前数据库事件,包括：insert、update、delete

local id = row["id"]

print(id)

if type(id)=="nil" then
	return
end

local result = {}  -- 定义一个table,作为结果集，其mappings在rule->es_mappings中定义

result["id"] = id
result["username"] = row["username"]
result["age"] = row["age"]
result["addr"] = row["addr"]

if action == "insert" then -- 新增事件
	ops.INSERT("artcool_product",id,result) -- 新增，参数1为index名称，string类型；参数2为要插入的数据主键；参数3为要插入的数据，tablele类型或者json字符串
elseif action == "delete" then -- 删除事件
	ops.DELETE("artcool_product",id) -- 删除，参数1为index名称，string类型；参数2为要插入的数据主键
else -- 修改事件
	ops.UPDATE("artcool_product",id,result) -- 修改，参数1为index名称，string类型；参数2为要插入的数据主键；参数3为要插入的数据，tablele类型或者json字符串
end