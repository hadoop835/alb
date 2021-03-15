local ngx = ngx
local ngx_header = ngx.header

local matched_policy = ngx.ctx.matched_policy
if matched_policy ~= nil then
    local enable_cors = matched_policy["enable_cors"]
    if enable_cors == true then
        if ngx.ctx.alb_ctx.method ~= 'OPTIONS' then
            if matched_policy["cors_allow_origin"] ~= "" then
                ngx_header['Access-Control-Allow-Origin'] = matched_policy["cors_allow_origin"]
            else
                ngx_header['Access-Control-Allow-Origin'] = "*"
            end
            ngx_header['Access-Control-Allow-Credentials'] = 'true'
            ngx_header['Access-Control-Allow-Methods'] = 'GET, PUT, POST, DELETE, PATCH, OPTIONS'
            if matched_policy["cors_allow_headers"] ~= "" then
                ngx_header['Access-Control-Allow-Headers'] = matched_policy["cors_allow_headers"]
            else
                ngx_header['Access-Control-Allow-Headers'] = 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization'
            end
        end
    end
end