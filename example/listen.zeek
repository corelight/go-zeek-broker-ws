redef exit_only_after_terminate = T;

global test_topic = "/topic/test";


export {
    global ping: event(msg: string, c: count);
    global pong: event(msg: string, c: count);
}

event ping(msg: string, c: count)
        {
        print fmt("receiver got ping: %s, %s", msg, c);
        local e = Broker::make_event(pong, msg, c+1);
        Broker::publish(test_topic, e);
        }

event zeek_init()
    {
    Broker::listen_websocket();
    Broker::subscribe(test_topic);
    }

event Broker::peer_added(endpoint: Broker::EndpointInfo, msg: string)
    {
    print "peer added", endpoint;
    }

event Broker::peer_lost(endpoint: Broker::EndpointInfo, msg: string)
    {
    print "peer lost", endpoint;
    terminate();
    }

