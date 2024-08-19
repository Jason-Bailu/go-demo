package redis_distributed_lock

// 判断是否拥有分布式锁的归属权，是则删除
const LuaCheckAndDeleteDistributedLock = `
	local lockerKey = KEYS[1]
	local targetToken = ARGV[1]
	local getToken = redis.call('get', lockerKey)
	if (not getToken or getToken ~= targetToken) then
		return 0
	else
		return redis.call('del', lockerKey)
	end
`

// 判断是否拥有分布式锁的归属权，然后续签
const LuaCheckAndExpiredDistributedLock = `
	local lockerKey = KEYS[1]
	local targetToken = ARGV[1]
	local duration = ARGV[2]
	local getToken = redis.call('get', lockerKey)
	if (not getToken or getToken ~= targetToken) then
		return 0
	else
		return redis.call('expire', lockerKey, duration)
	end
`
