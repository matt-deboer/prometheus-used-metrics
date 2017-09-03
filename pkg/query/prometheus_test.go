package query

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
)

const testRulesData = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>Prometheus Time Series Collection and Processing Server</title>
    <link rel="shortcut icon" href="/static/img/favicon.ico?v=3afb3fffa3a29c3de865e1172fb740442e9d0133">
    <script src="/static/vendor/js/jquery.min.js?v=3afb3fffa3a29c3de865e1172fb740442e9d0133"></script>
    <script src="/static/vendor/bootstrap-3.3.1/js/bootstrap.min.js?v=3afb3fffa3a29c3de865e1172fb740442e9d0133"></script>

    <link type="text/css" rel="stylesheet" href="/static/vendor/bootstrap-3.3.1/css/bootstrap.min.css?v=3afb3fffa3a29c3de865e1172fb740442e9d0133">
    <link type="text/css" rel="stylesheet" href="/static/css/prometheus.css?v=3afb3fffa3a29c3de865e1172fb740442e9d0133">

    <script>
      var PATH_PREFIX = "";
      var BUILD_VERSION = "3afb3fffa3a29c3de865e1172fb740442e9d0133";
      $(function () {
        $('[data-toggle="tooltip"]').tooltip()
      })
    </script>

    
  </head>

  <body>
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="/">Prometheus</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
          <ul class="nav navbar-nav navbar-left">
            
            
            <li><a href="/alerts">Alerts</a></li>
            <li><a href="/graph">Graph</a></li>
            <li class="dropdown">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">Status <span class="caret"></span></a>
              <ul class="dropdown-menu">
                <li><a href="/status">Runtime &amp; Build Information</a></li>
                <li><a href="/flags">Command-Line Flags</a></li>
                <li><a href="/config">Configuration</a></li>
                <li><a href="/rules">Rules</a></li>
                <li><a href="/targets">Targets</a></li>
              </ul>
            </li>
            <li>
              <a href="https://prometheus.io" target="_blank">Help</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    
  <div class="container-fluid">
    <h2 id="rules">Rules</h2>
    <pre>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_master_low_disk%22%7D&g0.tab=0">system_mesos_master_low_disk</a>
  IF <a href="/graph?g0.expr=%28node_filesystem_free%7Bmesos_role%3D%22master%22%2Cmountpoint%3D%22%2Frootfs%22%7D+%2F+node_filesystem_size%7Bmesos_role%3D%22master%22%2Cmountpoint%3D%22%2Frootfs%22%7D+%2A+100%29+%3C+20&g0.tab=0">(node_filesystem_free{mesos_role=&#34;master&#34;,mountpoint=&#34;/rootfs&#34;} / node_filesystem_size{mesos_role=&#34;master&#34;,mountpoint=&#34;/rootfs&#34;} * 100) &lt; 20</a>
  FOR 2h
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;One or more Mesos masters has under 20% disk available on /&#34;, summary=&#34;Mesos Master disk alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_agent_low_disk_root%22%7D&g0.tab=0">system_mesos_agent_low_disk_root</a>
  IF <a href="/graph?g0.expr=%28node_filesystem_free%7Bmesos_role%3D%22agent%22%2Cmountpoint%3D%22%2Frootfs%22%7D+%2F+node_filesystem_size%7Bmesos_role%3D%22agent%22%2Cmountpoint%3D%22%2Frootfs%22%7D+%2A+100%29+%3C+10&g0.tab=0">(node_filesystem_free{mesos_role=&#34;agent&#34;,mountpoint=&#34;/rootfs&#34;} / node_filesystem_size{mesos_role=&#34;agent&#34;,mountpoint=&#34;/rootfs&#34;} * 100) &lt; 10</a>
  FOR 2h
  LABELS {mountpoint=&#34;/&#34;, owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;One or more Mesos agents has under 10% disk available on /&#34;, summary=&#34;Mesos Agent disk alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_agent_low_disk_var_lib%22%7D&g0.tab=0">system_mesos_agent_low_disk_var_lib</a>
  IF <a href="/graph?g0.expr=%28node_filesystem_free%7Bmesos_role%3D%22agent%22%2Cmountpoint%3D%22%2Frootfs%2Fvar%2Flib%22%7D+%2F+node_filesystem_size%7Bmesos_role%3D%22agent%22%2Cmountpoint%3D%22%2Frootfs%2Fvar%2Flib%22%7D+%2A+100%29+%3C+10&g0.tab=0">(node_filesystem_free{mesos_role=&#34;agent&#34;,mountpoint=&#34;/rootfs/var/lib&#34;} / node_filesystem_size{mesos_role=&#34;agent&#34;,mountpoint=&#34;/rootfs/var/lib&#34;} * 100) &lt; 10</a>
  FOR 2h
  LABELS {mountpoint=&#34;/var/lib&#34;, owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;One or more Mesos agents has under 10% disk available on /var/lib&#34;, summary=&#34;Mesos Agent disk alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_marathonlb_internal_not_running%22%7D&g0.tab=0">system_mesos_marathonlb_internal_not_running</a>
  IF <a href="/graph?g0.expr=%28absent%28marathon_app_instances%7Bapp%3D%22%2Fsystem%2Fmarathon-lb-internal%22%7D%29+or+sum%28marathon_app_task_running%7Bapp%3D%22%2Fsystem%2Fmarathon-lb-internal%22%7D%29+%3C+3%29&g0.tab=0">(absent(marathon_app_instances{app=&#34;/system/marathon-lb-internal&#34;}) or sum(marathon_app_task_running{app=&#34;/system/marathon-lb-internal&#34;}) &lt; 3)</a>
  FOR 5m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;One or more INTERNAL ingress load balancer instances (out of 3) are not running.&#34;, summary=&#34;Mesos ingress load balancer alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_marathonlb_external_not_running%22%7D&g0.tab=0">system_mesos_marathonlb_external_not_running</a>
  IF <a href="/graph?g0.expr=%28absent%28marathon_app_instances%7Bapp%3D%22%2Fsystem%2Fmarathon-lb-external%22%7D%29+or+sum%28marathon_app_task_running%7Bapp%3D%22%2Fsystem%2Fmarathon-lb-external%22%7D%29+%3C+3%29&g0.tab=0">(absent(marathon_app_instances{app=&#34;/system/marathon-lb-external&#34;}) or sum(marathon_app_task_running{app=&#34;/system/marathon-lb-external&#34;}) &lt; 3)</a>
  FOR 5m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;One or more EXTERNAL ingress load balancer instances (out of 3) are not running.&#34;, summary=&#34;Mesos ingress load balancer alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22flapping__marathon_app%22%7D&g0.tab=0">flapping__marathon_app</a>
  IF <a href="/graph?g0.expr=marathon_app_task_running+%3E+0+and+max_over_time%28marathon_app_task_avg_uptime%5B10m%5D%29+%3C+30&g0.tab=0">marathon_app_task_running &gt; 0 and max_over_time(marathon_app_task_avg_uptime[10m]) &lt; 30</a>
  FOR 10m
  LABELS {owner=&#34;system&#34;, severity=&#34;warning&#34;}
  ANNOTATIONS {description=&#34;A Service (Marathon app) has been flapping for at least 10 minutes&#34;, summary=&#34;Marathon App/Service alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_flapping_marathon%22%7D&g0.tab=0">system_flapping_marathon</a>
  IF <a href="/graph?g0.expr=marathon_service_mesosphere_marathon_uptime+%3C+60&g0.tab=0">marathon_service_mesosphere_marathon_uptime &lt; 60</a>
  FOR 5m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;Marathon has been flapping for at least 10 minutes&#34;, summary=&#34;Marathon alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_flapping_mesos_master%22%7D&g0.tab=0">system_flapping_mesos_master</a>
  IF <a href="/graph?g0.expr=mesos_master_uptime_seconds+%3C+300&g0.tab=0">mesos_master_uptime_seconds &lt; 300</a>
  FOR 15m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;A Mesos master has been flapping for at least 15 minutes&#34;, summary=&#34;Mesos alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_flapping_mesos_agent%22%7D&g0.tab=0">system_flapping_mesos_agent</a>
  IF <a href="/graph?g0.expr=mesos_slave_uptime_seconds+%3C+300&g0.tab=0">mesos_slave_uptime_seconds &lt; 300</a>
  FOR 15m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;A Mesos agent has been flapping for at least 15 minutes&#34;, summary=&#34;Mesos alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22BlackboxExporterDown%22%7D&g0.tab=0">BlackboxExporterDown</a>
  IF <a href="/graph?g0.expr=count%28up%7Bjob%3D%22blackbox%22%7D+%3D%3D+1%29+%3C+1&g0.tab=0">count(up{job=&#34;blackbox&#34;} == 1) &lt; 1</a>
  FOR 3m
  LABELS {ci=&#34;DCOS&#34;, noc_alert=&#34;f&#34;, owner=&#34;system&#34;, selfservice=&#34;t&#34;, severity=&#34;critical&#34;, team=&#34;Ops&#34;, tsa=&#34;KB0013754&#34;}
  ANNOTATIONS {description=&#34;The blackbox exporter (which provides other critical alerts) has been down for at least 3 minutes&#34;, summary=&#34;Blackbox exporter alert&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22container_failure_threshold_exceeded%22%7D&g0.tab=0">container_failure_threshold_exceeded</a>
  IF <a href="/graph?g0.expr=%28sum%28increase%28mesos_slave_task_states_exit_total%7Bstate%21~%22killed%7Cfinished%22%7D%5B5m%5D%29%29+%2F+sum%28mesos_slave_task_states_current%29+%3E+0.9%29&g0.tab=0">(sum(increase(mesos_slave_task_states_exit_total{state!~&#34;killed|finished&#34;}[5m])) / sum(mesos_slave_task_states_current) &gt; 0.9)</a>
  FOR 5m
  LABELS {owner=&#34;system&#34;, severity=&#34;warning&#34;}
  ANNOTATIONS {description=&#34;Too many containers are failing&#34;, summary=&#34;The Container Failure Rate is too damn high&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_mesos_critical_url_down%22%7D&g0.tab=0">system_mesos_critical_url_down</a>
  IF <a href="/graph?g0.expr=%28probe_http_status_code+%21%3D+200%29&g0.tab=0">(probe_http_status_code != 200)</a>
  FOR 1m
  LABELS {ci=&#34;DCOS&#34;, noc_alert=&#34;f&#34;, owner=&#34;system&#34;, selfservice=&#34;t&#34;, severity=&#34;critical&#34;, team=&#34;Ops&#34;, tsa=&#34;KB0013754&#34;}
  ANNOTATIONS {description=&#34;We received a non-200 response from {{ $labels.instance }} for 5 minutes&#34;, summary=&#34;URL is down for {{ $labels.instance }}&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_low_cpu_in_rack%22%7D&g0.tab=0">system_low_cpu_in_rack</a>
  IF <a href="/graph?g0.expr=sum%28mesos_slave_cpus%7Bingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Crack%21%3D%22%22%2Ctype%3D%22used%22%7D%29+BY+%28rack%29+%3E+%28sum%28mesos_slave_cpus%7Bingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Crack%21%3D%22%22%7D%29+BY+%28rack%29+%2A+0.9%29&g0.tab=0">sum(mesos_slave_cpus{ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,rack!=&#34;&#34;,type=&#34;used&#34;}) BY (rack) &gt; (sum(mesos_slave_cpus{ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,rack!=&#34;&#34;}) BY (rack) * 0.9)</a>
  FOR 30m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;Less than 10% CPU available in rack {{ $labels.rack }}&#34;, summary=&#34;Approaching maximum CPU usage in rack {{ $labels.rack }}&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_low_cpu_in_az%22%7D&g0.tab=0">system_low_cpu_in_az</a>
  IF <a href="/graph?g0.expr=sum%28mesos_slave_cpus%7Baz%21%3D%22%22%2Cingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Ctype%3D%22used%22%7D%29+BY+%28az%29+%3E+%28sum%28mesos_slave_cpus%7Baz%21%3D%22%22%2Cingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%7D%29+BY+%28az%29+%2A+0.9%29&g0.tab=0">sum(mesos_slave_cpus{az!=&#34;&#34;,ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,type=&#34;used&#34;}) BY (az) &gt; (sum(mesos_slave_cpus{az!=&#34;&#34;,ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;}) BY (az) * 0.9)</a>
  FOR 30m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;Less than 10% CPU available in az {{ $labels.az }}&#34;, summary=&#34;Approaching maximum CPU usage in az {{ $labels.az }}&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_low_mem_in_rack%22%7D&g0.tab=0">system_low_mem_in_rack</a>
  IF <a href="/graph?g0.expr=sum%28mesos_slave_mem%7Bingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Crack%21%3D%22%22%2Ctype%3D%22used%22%7D%29+BY+%28rack%29+%3E+%28sum%28mesos_slave_mem%7Bingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Crack%21%3D%22%22%7D%29+BY+%28rack%29+%2A+0.9%29&g0.tab=0">sum(mesos_slave_mem{ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,rack!=&#34;&#34;,type=&#34;used&#34;}) BY (rack) &gt; (sum(mesos_slave_mem{ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,rack!=&#34;&#34;}) BY (rack) * 0.9)</a>
  FOR 30m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;Less than 10% memory available in rack {{ $labels.rack }}&#34;, summary=&#34;Approaching maximum memory usage in rack {{ $labels.rack }}&#34;}<br/>ALERT <a href="/graph?g0.expr=ALERTS%7Balertname%3D%22system_low_mem_in_az%22%7D&g0.tab=0">system_low_mem_in_az</a>
  IF <a href="/graph?g0.expr=sum%28mesos_slave_mem%7Baz%21%3D%22%22%2Cingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%2Ctype%3D%22used%22%7D%29+BY+%28az%29+%3E+%28sum%28mesos_slave_mem%7Baz%21%3D%22%22%2Cingress%3D%22%22%2Cjob%3D%22mesos_exporter_agents%22%7D%29+BY+%28az%29+%2A+0.9%29&g0.tab=0">sum(mesos_slave_mem{az!=&#34;&#34;,ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;,type=&#34;used&#34;}) BY (az) &gt; (sum(mesos_slave_mem{az!=&#34;&#34;,ingress=&#34;&#34;,job=&#34;mesos_exporter_agents&#34;}) BY (az) * 0.9)</a>
  FOR 30m
  LABELS {owner=&#34;system&#34;, severity=&#34;critical&#34;}
  ANNOTATIONS {description=&#34;Less than 10% memory available in az {{ $labels.az }}&#34;, summary=&#34;Approaching maximum memory usage in az {{ $labels.az }}&#34;}<br/></pre>
  </div>

  </body>
</html>

`

type mockPrometheus struct {
}

func (mp *mockPrometheus) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/rules" {
		w.Write([]byte(testRulesData))
	} else {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

}

func TestGetAlertMetricUsage(t *testing.T) {

	prom1 := &mockPrometheus{}
	prom1Server := httptest.NewServer(prom1)
	defer prom1Server.Close()

	pr, err := NewPrometheusMetricsResolver(prom1Server.URL)
	if err != nil {
		t.Error(err)
	}

	log.SetLevel(log.DebugLevel)

	usage, err := pr.GetAlertMetricUsage()
	if err != nil {
		t.Error(err)
	}

	if usage == nil {
		t.Errorf("Expected usage metrics")
	}
	if len(usage) != 14 {
		t.Errorf("Expected 14 metrics used, got %d", len(usage))
	}
}
