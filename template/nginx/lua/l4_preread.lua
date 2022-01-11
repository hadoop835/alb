local ngx = ngx
local ngx_var = ngx.var
local ngx_log = ngx.log
local ngx_exit = ngx.exit

local upstream = require "upstream"
local var_proxy = require "var_proxy"

local subsystem = ngx.config.subsystem
ngx.ctx.alb_ctx = var_proxy.new()

local t_upstream, _, errmsg = upstream.get_upstream(subsystem, ngx_var.protocol, ngx_var.server_port)
if t_upstream ~= nil then
    ngx_var.upstream = t_upstream
end

if errmsg ~= nil then
    ngx_log(ngx.ERR, errmsg)
    ngx_exit(ngx.ERROR)
end
