package dllm.app;

import android.os.Bundle;
import com.getcapacitor.BridgeActivity;
import dllm.app.plugins.Client.ClientPlugin;

public class MainActivity extends BridgeActivity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        registerPlugin(ClientPlugin.class);
        super.onCreate(savedInstanceState);
    }
}
