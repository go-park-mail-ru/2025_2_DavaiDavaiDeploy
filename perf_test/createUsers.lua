wrk.method = "POST"
local counter = 0
local total_users = 100000

local unique_prefix = tostring(math.random(1000000000, 9999999999))
math.randomseed(os.time() * 1000)

request = function()
    counter = counter + 1
    if counter > total_users then
        return nil
    end
    
    local login = string.format("%s_%d_%d", 
        unique_prefix, 
        counter, 
        math.random(1000, 9999))
    

    login = login:sub(1, 20)
    
    if #login < 6 then
        login = login .. string.rep("x", 6 - #login)
    end
    
    local body = string.format('{"login": "%s", "password": "password123"}', login)
    
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    
    return wrk.format("POST", "/api/auth/signup", headers, body)
end

done = function(summary, latency, requests)
    print(string.format("\nCreated %d users", counter - 1))
end