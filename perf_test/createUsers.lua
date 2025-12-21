wrk.method = "POST"
local counter = 0
local total_users = 100000
local password = "pass123"

local date_str = os.date("%m%d%H%M")  
local random_suffix = math.random(100, 999) 
local prefix = "u" .. date_str .. random_suffix  
math.randomseed(os.time() * 1000)

request = function()
    counter = counter + 1
    if counter > total_users then
        return nil
    end
    local random_part = math.random(100, 999)
    local login = string.format("%s_%d_%d", prefix, counter, random_part)
    
    login = login:sub(1, 15)
    
    if #login < 6 then
        login = login .. string.rep("a", 6 - #login)
    end
    
    local body = string.format('{"login": "%s", "password": "%s"}', login, password)
    
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    
    return wrk.format("POST", "/api/auth/signup", headers, body)
end

done = function(summary, latency, requests)
    print(string.format("\nCreated %d users with prefix: %s", counter, prefix))
end