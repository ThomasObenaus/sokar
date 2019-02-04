job "fail-service" {
  datacenters = ["public-services"]

  type = "service"

  update {
    stagger = "5s"
    max_parallel = 1
  }
      

  group "fail-service-grp-A" {
    count = 1
    task "fail-service" {
      driver = "docker"
      config {
        image = "thobe/fail_service:latest"
        port_map = {
          http = 8080
        }
      }

      # Register at consul
      service {
        name = "${TASK}"
        port = "http"
        check {
          port     = "http"
          type     = "http"
          path     = "/health"
          method   = "GET"
          interval = "10s"
          timeout  = "2s"
        }
      }

      env {
        HEALTHY_FOR   = -1,
      }

      resources {
        cpu    = 100 # MHz
        memory = 256 # MB
        network {
          mbits = 10
          port "http" {}
        }
      }
    }
  }


  group "fail-service-grp-B" {
    count = 2
    task "fail-service" {
      driver = "docker"
      config {
        image = "thobe/fail_service:latest"
        port_map = {
          http = 8080
        }
      }

      # Register at consul
      service {
        name = "${TASK}"
        port = "http"
        check {
          port     = "http"
          type     = "http"
          path     = "/health"
          method   = "GET"
          interval = "10s"
          timeout  = "2s"
        }
      }

      env {
        HEALTHY_FOR   = -1,
      }

      resources {
        cpu    = 100 # MHz
        memory = 256 # MB
        network {
          mbits = 10
          port "http" {}
        }
      }
    }
  }
}
