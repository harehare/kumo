using System;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Kumo.Models {
  public class Item {
    [JsonPropertyName ("user_id")]
    public string UserID { get; set; }

    [Required]
    [JsonPropertyName ("name")]
    [StringLength (100)]
    public string Name { get; set; }

    [JsonPropertyName ("description")]
    public string Description { get; set; }

    [Url]
    [Required]
    [JsonPropertyName ("url")]
    public string Url { get; set; }

    [Required]
    [JsonPropertyName ("selector")]
    public string Selector { get; set; }

    [Required]
    [ValidateComplexType]
    [JsonPropertyName ("notification")]
    public Notification Notification { get; set; }

    [JsonPropertyName ("time")]
    public DateTime? CreatedAt { get; set; }

    public Item () {
      Name = "";
      Description = "";
      Url = "";
      Selector = "";
      Notification = new Notification ();
      CreatedAt = null;
    }

    public Item (string url, string selector) {
      this.Url = url;
      this.Selector = selector;
    }

    public override bool Equals (Object obj) {
      if (obj == null || obj.GetType () != GetType ()) {
        return false;
      }
      var item = obj as Item;
      return item.Selector == Selector && item.Url == Url;
    }
    public override int GetHashCode () {
      return Selector.GetHashCode () ^ Url.GetHashCode ();
    }
  }

  public class History {
    public string ID { get; set; }
    public string Status { get; set; }
    public string Text { get; set; }
    public DateTime Timestamp { get; set; }
  }
}
