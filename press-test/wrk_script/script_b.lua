wrk.headers["content-type"] = "application/json"
local rawdata_path = "../rawdata"
local file_index = 1
local event_index = 1
-- local events = {}
-- local total_event_count = 0
-- local request_start_time
local event
local log_file

function init(args)
  endpoint = args[1]
  file_list = get_file_list(rawdata_path)
  total_file_count = table.getn(file_list)
  log_file = io.open(string.format("../logs/wrk_script/output_%s.log", now()), "w")
end

local ffi = require "ffi"

ffi.cdef [[
    typedef long time_t;
    typedef long suseconds_t;
    typedef struct timeval {
               time_t      tv_sec;     
               suseconds_t tv_usec;    
           } timeval_t;
    typedef struct timezone {
               int tz_minuteswest;     
               int tz_dsttime;         
           }timezone_t;
    int gettimeofday(struct timeval *tv, struct timezone *tz);
]]

function now()
    local timeval_t = ffi.typeof("timeval_t")
    local tv = ffi.new(timeval_t)
    local timezone_t = ffi.typeof("timezone_t")
    local tz = ffi.new(timezone_t)
    ffi.C.gettimeofday(tv, tz)
    return tonumber(tv.tv_sec) + tonumber(tv.tv_usec) / (1000 * 1000)
end

-- Get the file names from rawdata path
function get_file_list(path)
  local a = io.popen("ls "..path);
  local f = {};
  for l in a:lines() do
      table.insert(f,l)
  end
  a:close()
  return f
end

-- -- Get the events from a specific file
-- function get_events()
--   local file = string.format("%s/%s", rawdata_path, file_list[file_index])
--   local start_time = now();
--   local lines = {}
--   for line in io.lines(file) do 
--     lines[#lines + 1] = line
--   end
--   local end_time = now();
--   log_file:write(string.format("Get events from file %s start %s end %s took %s s \n", file, start_time, end_time, end_time - start_time))
--   return lines
-- end

-- Get the events from a specific file
function get_event_from_file(f_index, e_index)
  local file = string.format("%s/%s", rawdata_path, file_list[f_index])
  local start_time = now();
  local i = 1
  local e
  for line in io.lines(file) do 
    if i == event_index then
      e = line
      event_index = event_index + 1
      break
    else
      i = i + 1
    end
  end
  local end_time = now();
  log_file:write(string.format("Get #%s event from file %s start %s end %s took %s s \n", e_index, file, start_time, end_time, end_time - start_time))
  return e
end

-- Get an event
function get_event()
  event = get_event_from_file(file_index, event_index)
  if event == nil or event == "" then
    file_index = file_index + 1
    event_index = 1
    event = get_event_from_file(file_index, event_index)
  end
  return event
end

-- -- Get an event
-- function get_event()
--   if event_index < total_event_count then
--     event_index = event_index + 1
--   else
--     -- stop wrk if the index in the last file and the last line
--     if total_file_count == file_index then
--       log_file:write("No more data...")
--       log_file:close()
--       wrk.thread:stop()
--       return
--     end

--     -- Switch to next file
--     file_index = file_index + 1
--     event_index = 1
--     events = get_events()
--     total_event_count = table.getn(events)
--   end
--   return events[event_index]
-- end

function request()
  request_start_time = now()
  event = get_event()
  return wrk.format("POST", endpoint, nil, event)
end

function response(status, headers, body)
  local cost_time = now() - request_start_time
  if cost_time > 0.5 then
    log_file:write(string.format("Event: %s\nStatus: %s\nTook: %s s\nResponse: %s\n", event, status, cost_time, body))
  end
  -- print(cjson.encode(headers))
  -- print(status)
  -- print(body)
end

