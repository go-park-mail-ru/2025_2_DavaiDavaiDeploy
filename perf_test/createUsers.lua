wrk.method = "POST"
local counter = 100000 
local total_users = 100000
local password = "Test123!"

request = function()
    counter = counter + 1
    if counter > 200000 then 
        return nil
    end
    
    local login = "user" .. tostring(counter)
    
    
    local body = string.format('{"login": "%s", "password": "%s"}', login, password)
    
    local headers = {
        ["Content-Type"] = "application/json",
        ["Accept"] = "application/json"
    }
    
    return wrk.format("POST", "/api/auth/signup", headers, body)
end

done = function(summary, latency, requests)
    print(string.format("\nCreated %d users", counter - 100000))
end