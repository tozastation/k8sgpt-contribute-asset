resource "google_pubsub_topic" "gke_cluster_notification_topic" {
  name = "${var.resource_prefix}-gke-cluster-notification-topic"
  # 7days
  message_retention_duration = "1512000s"
}

resource "google_pubsub_subscription" "gke_cluster_notification_pull_subscription" {
  name  = "${var.resource_prefix}-gke-cluster-notification-pull-subscription"
  topic = google_pubsub_topic.gke_cluster_notification_topic.id

  # 7days
  message_retention_duration   = "604800s"
  enable_exactly_once_delivery = true
}