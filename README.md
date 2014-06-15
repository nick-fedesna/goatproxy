GoatProxy
=========

A small HTTP proxy written in Go with option to launch an Android Activity passing necessary proxy info as extras.

## Command Usage

        $ ./goatproxy -host=my.app.com -port:8888 -ssl=true -latency=0 -pkg=com.example.app

## Android
Requires permission `android.permission.ACCESS_WIFI_STATE`.

Drop the code below into your project, call `GoatProxy.getGoatProxy(this)` in onCreate() of your main launch activity which returns the `Proxy` or `Proxy.NO_PROXY`.   

SSL: override your `https://` network requests to use `http://` when using the Proxy.
    
    import static android.content.Context.WIFI_SERVICE;
    import static android.util.Patterns.IP_ADDRESS;
    
    import android.app.Activity;
    import android.net.DhcpInfo;
    import android.net.wifi.WifiManager;
    
    import java.net.*;
    
    public class GoatProxy {
    
        public static Proxy getGoatProxy(Activity aActivity) {
            if (aActivity.getIntent().hasExtra("goatProxyHosts")) {
                String[] ips = aActivity.getIntent().getStringExtra("goatProxyHosts").split("�");
                if (ips != null && ips.length > 0) {
                    WifiManager wifi = (WifiManager) aActivity.getSystemService(WIFI_SERVICE);
                    int myAddr = wifi.getConnectionInfo().getIpAddress();
                    InetAddress ip = matchProxyAddress(myAddr, ips, wifi.getDhcpInfo());
                    if (ip != null) {
                        String port = aActivity.getIntent().getStringExtra("goatProxyPort");
                        port = (port == null) ? "8080" : port;
                        return new Proxy(Proxy.Type.HTTP, new InetSocketAddress(ip, Integer.parseInt(port)));
                    }
                }
            }
            return Proxy.NO_PROXY;
        }
    
        private static InetAddress matchProxyAddress(int aIpAddr, String[] aIpArray, DhcpInfo aDhcp) {
            for (String ip : aIpArray) {
                if (IP_ADDRESS.matcher(ip).matches()) {
                    try {
                        InetAddress[] inets = InetAddress.getAllByName(ip);
                        if (inets != null && inets.length > 0) {
                            int sub = inetAddressToInt(inets[0]) & 0xffffff;
                            if (sub == (aIpAddr & 0xffffff)) {
                                return inets[0];
                            }
                            if (aDhcp != null && sub == (aDhcp.dns1 & 0xffffff)) {
                                return inets[0];
                            }
                        }
                    } catch (UnknownHostException ignored) { }
                }
            }
    
            return null;
        }
    
        private static int inetAddressToInt(InetAddress inetAddr) throws IllegalArgumentException {
            byte[] addr = inetAddr.getAddress();
            return ((addr[3] & 0xff) << 24) | ((addr[2] & 0xff) << 16) |
                    ((addr[1] & 0xff) << 8) | (addr[0] & 0xff);
        }
    
    }

	
###To-Do:
* Android: remember Proxy info and detect if it is still active (when extras not present)
* Android: hook into Activity lifecycle callbacks to get launch intent extras, removing need to pass Activity to getGoatProxy()
* Go: present device chooser when multiple adb devices found