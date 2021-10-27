local key = KEYS[1]
local ttlMs = ARGV[1]

if not redis.call("GET", key) then
    -- the key is absent in the storage
    redis.call("SET", key, 1)
    redis.call("PEXPIRE", key, ttlMs)
    return 1
end

return redis.call("INCR", key)