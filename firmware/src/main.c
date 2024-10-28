// // Camera
#include "camera/camera.h"
#include "esp_camera.h"

// // Wifi
#include "wifi/wifi.h"

// // Http server
#include "http/http.h"
#include "esp_http_server.h"

// Extras
#include "esp_log.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

static const char *TAG = "MAIN";

// Handler function for HTTP GET query
esp_err_t get_stream_handler(httpd_req_t *req)
{
    char *part_buf[64];
    while (true) {
        camera_fb_t *pic = esp_camera_fb_get();
        if (!pic) {
            ESP_LOGE(TAG, "Failed to get camera frame buffer");
            httpd_resp_send_500(req);
            return ESP_FAIL;
        }

        // Send multipart header
        size_t hlen = snprintf((char *)part_buf, 64, "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %u\r\n\r\n", pic->len);
        if (httpd_resp_send_chunk(req, (const char *)part_buf, hlen) != ESP_OK) {
            esp_camera_fb_return(pic);
            break;
        }

        // Send frame
        if (httpd_resp_send_chunk(req, (const char *)pic->buf, pic->len) != ESP_OK) {
            esp_camera_fb_return(pic);
            break;
        }

        // End of frame
        if (httpd_resp_send_chunk(req, "\r\n", 2) != ESP_OK) {
            esp_camera_fb_return(pic);
            break;
        }

        esp_camera_fb_return(pic);
        vTaskDelay(pdMS_TO_TICKS(100));  // Adjust delay for frame rate
    }
    return ESP_OK;
}

httpd_uri_t uri_get_stream = {
    .uri      = "/stream",
    .method   = HTTP_GET,
    .handler  = get_stream_handler,
    .user_ctx = NULL
};

httpd_uri_t uri_get_capture = {
    .uri      = "/capture",
    .method   = HTTP_GET,
    .handler  = get_stream_handler,
    .user_ctx = NULL
};


void app_main() {
    ESP_ERROR_CHECK(setup_wifi());
    ESP_ERROR_CHECK(init_camera());

    httpd_handle_t webServer = start_webserver();
    if (webServer == NULL) {
        ESP_LOGE(TAG, "Failed to start web server");
        return;
    }
    ESP_LOGI(TAG, "Web server started");

    httpd_register_uri_handler(webServer, &uri_get_stream);
    httpd_register_uri_handler(webServer, &uri_get_capture);

    while(1) {
        ESP_LOGI(TAG, "WAITING");
        vTaskDelay(pdMS_TO_TICKS(1000));
    }
}