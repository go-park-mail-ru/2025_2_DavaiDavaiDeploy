wrk.method = "POST"

local counter = 1
local max_users = 1000 
local prefix = "user_" .. os.time() .. "_"


math.randomseed(os.time())

request = function()
    if counter > max_users then
        return nil  
    end
    
    local login = prefix .. counter .. "_" .. math.random(10000, 99999)
    local password = "password123"
    
    local body = string.format('{"login": "%s", "password": "%s"}', login, password)
    
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    
    counter = counter + 1
    return wrk.format("POST", "/api/auth/signup", headers, body)
end

done = function(summary, latency, requests)
    print("\nСоздано пользователей: " .. (counter - 1) ..)
end
