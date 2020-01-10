local ngx_var = ngx.var
local ngx_header = ngx.header
local ngx_re = ngx.re
local ngx_req = ngx.req
local ngx_log = ngx.log
local ngx_say = ngx.say
local ngx_exit = ngx.exit
local upstream = require "upstream"

local t_upstream, matched_policy, errmsg = upstream.get_upstream(ngx_var.server_port)
if t_upstream ~= nil then
  ngx_var.upstream = t_upstream
end
if matched_policy ~= nil then
  ngx_var.rule_name = matched_policy["rule"]
  local enable_cors = matched_policy["enable_cors"]
  if enable_cors == true then
    if ngx_req.get_method() == 'OPTION' then
      ngx_header['Access-Control-Allow-Origin'] = '*'
      ngx_header['Access-Control-Allow-Credentials'] = 'true'
      ngx_header['Access-Control-Allow-Methods'] = 'GET, PUT, POST, DELETE, PATCH, OPTIONS'
      ngx_header['Access-Control-Allow-Headers'] = 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization'
      ngx_header['Access-Control-Max-Age'] = '1728000'
      ngx_header['Content-Type'] = 'text/plain charset=UTF-8'
      ngx_header['Content-Length'] = '0'
      ngx_exit(ngx.HTTP_NO_CONTENT)
    else
      ngx_header['Access-Control-Allow-Origin']= '*'
      ngx_header['Access-Control-Allow-Credentials'] = 'true'
      ngx_header['Access-Control-Allow-Methods'] = 'GET, PUT, POST, DELETE, PATCH, OPTIONS'
      ngx_header['Access-Control-Allow-Headers'] = 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization'
    end
  end
  local backend_protocol = matched_policy["backend_protocol"]
  if backend_protocol ~= "" then
      ngx_var.backend_protocol = backend_protocol
  end
  local rewrite_target = matched_policy["rewrite_target"]
  local policy_url = matched_policy["url"]
  if rewrite_target ~= "" then
    if policy_url == "" then
      policy_url = "/"
    end
    local new_uri = ngx_re.sub(ngx_var.uri, policy_url, rewrite_target, "jo")
    ngx_req.set_uri(new_uri, false)
  end
elseif errmsg ~= nil then
  ngx.status = 404
  ngx_log(ngx.ERR, errmsg)
  ngx_say(errmsg)
  ngx_exit(ngx.HTTP_OK)
end
