wrk.method = "GET"

local user_ids = {}
local user_count = 0
local current_index = 1

function init(args)
    local file = io.open("user_ids.txt", "r")
    if not file then
        print("ERROR: Файл user_ids.txt не найден")
        os.exit(1)
    end
    
    print("Загружаем ID пользователей из файла...")
    for line in file:lines() do
        table.insert(user_ids, line)
    end
    file:close()
    
    user_count = #user_ids
    if user_count == 0 then
        print("ERROR: В файле нет ID пользователей!")
        os.exit(1)
    end
    
    print("Загружено " .. user_count .. " ID пользователей")
end

request = function()
    local user_id = user_ids[current_index]
    current_index = current_index + 1
    if current_index > user_count then
        current_index = 1
    end
    
    local path = "/api/users/" .. user_id
    
    local headers = {
        ["Accept"] = "application/json"
    }

    return wrk.format("GET", path, headers)
end

done = function(summary, latency, requests)
    print("\nТест чтения завершен")
    print("Использовано уникальных ID: " .. math.min(requests, user_count))
end