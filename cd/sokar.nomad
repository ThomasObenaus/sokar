# https://www.nomadproject.io/docs/job-specification/job.html
job "sokar" {
  datacenters = ["public-services"]

  type = "service"

  # https://www.nomadproject.io/docs/job-specification/reschedule.html
  reschedule {
    delay          = "30s"
    delay_function = "constant"
    unlimited      = true
  }

  # https://www.nomadproject.io/docs/job-specification/update.html
  update {
    max_parallel      = 1
    health_check      = "checks"
    min_healthy_time  = "10s"
    healthy_deadline  = "5m"
    progress_deadline = "10m"
    auto_revert       = true
    canary            = 0
    stagger           = "30s"
  }

  # https://www.nomadproject.io/docs/job-specification/group.html
  group "sokar" {
    count = 1

    # https://www.nomadproject.io/docs/job-specification/restart.html
    restart {
      interval = "10m"
      attempts = 2
      delay    = "15s"
      mode     = "fail"
    }

    # https://www.nomadproject.io/docs/job-specification/task.html
    task "sokar" {
      driver = "docker"

      config {
        image = "thobe/sokar:latest"
        port_map = {
          http = 11000
        }
        args = [
          "",
        ]
      }

      # https://www.nomadproject.io/docs/job-specification/env.html
      env {
        SK_NOMAD_SERVER_ADDRESS="http://${attr.unique.network.ip-address}:4646"
        SK_LOGGING_NO_COLOR="true"
      }

      # https://www.nomadproject.io/docs/job-specification/service.html
      service {
        name = "sokar"
        port = "http"
        tags = ["urlprefix-/sokar strip=/sokar"] # fabio

        check {
          name     = "sokar health using http endpoint '/health'"
          port     = "http"
          type     = "http"
          path     = "/health"
          method   = "GET"
          interval = "10s"
          timeout  = "2s"
        }

        # https://www.nomadproject.io/docs/job-specification/check_restart.html
        check_restart {
          limit = 3
          grace = "10s"
          ignore_warnings = false
        }
      }

      # https://www.nomadproject.io/docs/job-specification/resources.html
      resources {
        cpu    = 100 # MHz
        memory = 256 # MB

        # https://www.nomadproject.io/docs/job-specification/network.html
        network {
          mbits = 10
          port "http" {}
        }
      }
    }
  }
}