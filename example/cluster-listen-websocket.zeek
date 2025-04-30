redef exit_only_after_terminate = T;

global test_topic = "/topic/test";


export {
    global ping: event(msg: string, c: count);
    global pong: event(msg: string, c: count);
}

event ping(msg: string, c: count)
        {
        print fmt("receiver got ping: %s, %s", msg, c);
        local e = Cluster::make_event(pong, msg, c+1);
        Cluster::publish(test_topic, e);
        }

event zeek_init()
    {
    Cluster::listen_websocket([$listen_host="127.0.0.1", $listen_port=9997/tcp]);
    Cluster::subscribe(test_topic);
    }

event Cluster::websocket_client_added(endpoint: Cluster::EndpointInfo, subscriptions: string_vec)
    {
    print "web_socket_client_added", endpoint, subscriptions;
    }

event Cluster::websocket_client_lost(endpoint: Cluster::EndpointInfo)
    {
    print "websocket_client_lost", endpoint;
    terminate();
    }

