-- 获取缓存key
local key = KEYS[1]
-- 获取设置的值
local cntKey = ARGV[1]
--限制总数溢出情况
--tonumber这个函数会尝试将它的参数转换为数字
local delta = tonumber(ARGV[2])
-- 获取缓存已经存在的值
local exists = redis.call("EXISTS",key)
if exists == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    -- 说明自增成功了
    return 1
else
    return 0
end