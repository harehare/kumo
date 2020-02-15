using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Kumo.Models {
  public class Notification {
    [Url]
    [Required]
    [JsonPropertyName ("webhook_url")]
    public string WebHookUrl { get; set; }

    public Notification () {
      WebHookUrl = "";
    }
  }
}
