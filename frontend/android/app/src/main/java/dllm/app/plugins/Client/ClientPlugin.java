package dllm.app.plugins.Client;

import com.getcapacitor.Plugin;
import com.getcapacitor.PluginCall;
import com.getcapacitor.PluginMethod;
import com.getcapacitor.annotation.CapacitorPlugin;


@CapacitorPlugin(name = "Client")
public class ClientPlugin extends Plugin {
    @PluginMethod(returnType = PluginMethod.RETURN_CALLBACK)
    public void send(PluginCall call) {
        System.out.println("Sending message to server");
        call.setKeepAlive(true);

    }
}
