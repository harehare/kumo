using System.Text.Json.Serialization;

namespace Kumo.Models {
  public class User {
    [JsonPropertyName ("id")]
    public string ID { get; set; }

    [JsonPropertyName ("displayName")]
    public string DisplayName { get; set; }

    [JsonPropertyName ("photoUrl")]
    public string PhotoUrl { get; set; }

    [JsonPropertyName ("idToken")]
    public string IdToken { get; set; }
  }
}
