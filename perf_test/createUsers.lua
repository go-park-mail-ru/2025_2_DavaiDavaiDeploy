wrk.method = "POST"
local counter = 0
local total_users = 100000

function random_string(len)
    local chars = "abcdefghijklmnopqrstuvwxyz"
    local res = {}
    for i = 1, len do
        res[i] = chars:sub(math.random(1, #chars), math.random(1, #chars))
    end
    return table.concat(res)
end

request = function()
    counter = counter + 1
    if counter > total_users then
        return nil
    end
    
    local login = "user_" .. counter .. "_" .. random_string(6)
    
    local body = string.format('{"login": "%s", "password": "password123"}', login)
    
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    
    return wrk.format("POST", "/api/auth/signup", headers, body)
end

done = function(summary, latency, requests)
    print(string.format("\n✓ Created %d users", counter - 1))
end