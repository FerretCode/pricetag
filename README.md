# pricetag

A self hosted, simple, and robust log filtering platform with support for log forwarding for Railway

# features

-   deployed right in your Railway service
    1. deploy the Railway template
    2. expose your service via a domain or temporarily proxy it
    3. create an admin account for the environment
-   create log filters via "tags"
    -   filter by
        -   keyword
        -   JSON attribute
        -   service ID
-   log forwarding
    -   create pipelines for sending logs to other services via webhooks
    -   robust customization
        1. you're given variable names such as $LOG_CONTENT, $LOG_TIMESTAMP, etc.
        2. create the JSON object to be sent using the variable names
        3. point the forwarder at a webhook server
-   user access control
    -   view logs
        -   manage tags
        -   manage tracked services
        -   manage log forwarding
