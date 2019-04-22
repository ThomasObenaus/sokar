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
        # Needed to overwrite the --help (default argument of the sokar docker container)
        args = [
          "",
        ]
      }

      # https://www.nomadproject.io/docs/job-specification/env.html
      env {
        SK_NOMAD_SERVER_ADDRESS="http://${attr.unique.network.ip-address}:4646"
        SK_LOGGING_NO_COLOR="true"

        # This is just an example to show, how the needed job-configuration
        # for the fail-service job could look like.
        # You have to adjust all values according to your needs.
        SK_JOB_NAME="fail-service"
        SK_JOB_MIN=1
        SK_JOB_MAX=10
        SK_SAA_SCALE_ALERTS="AlertA:1.0:A upscaling alert;AlertB:-1.0:A downscaling alert"
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