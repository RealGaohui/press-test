wrk.headers["content-type"] = "application/json"
local rawdata_path = "../rawdata"
local file_index = 1
local event
local log_file
local file

function init(args)
  endpoint = args[1]
  file_list = get_file_list(rawdata_path)
  file = io.open(string.format("%s/%s", rawdata_path, file_list[file_index]), "r")
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

-- Get an event
function get_event()
  event = file:read()
  if event == nil or event == "" then
    if file_index == total_file_count then
      log_file:write("No more data...")
      log_file:close()
      file:close()
      wrk.thread:stop()
      return
    end
    file_index = file_index + 1
    file = io.open(string.format("%s/%s", rawdata_path, file_list[file_index]), "r")
    event = file:read()
  end
  return event
end

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

