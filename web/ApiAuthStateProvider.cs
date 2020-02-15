using System;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Security.Claims;
using System.Text.Json;
using System.Threading.Tasks;
using Kumo.Models;
using Microsoft.AspNetCore.Components.Authorization;
using Microsoft.JSInterop;

namespace Kumo.Auth {
  public class ApiAuthenticationStateProvider : AuthenticationStateProvider {
    private static string idToken = null;
    private static ClaimsPrincipal claimsPrincipal = new ClaimsPrincipal (new ClaimsIdentity ());
    private static ApiAuthenticationStateProvider provider = null;
    private readonly HttpClient httpClient;

    public ApiAuthenticationStateProvider (HttpClient httpClient) : base () {
      provider = this;
      this.httpClient = httpClient;
    }

    public override Task<AuthenticationState> GetAuthenticationStateAsync () {
      var authState = Task.FromResult (new AuthenticationState (claimsPrincipal));
      httpClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue ("Bearer", idToken);
      return authState;
    }

    public void NotifyAuthenticationStateChanged () {
      NotifyAuthenticationStateChanged (GetAuthenticationStateAsync ());
    }

    [JSInvokable]
    public static void MarkUserAsAuthenticated (string json) {
      var options = new JsonSerializerOptions { };
      var user = JsonSerializer.Deserialize<User> (json, options);
      var identity = new ClaimsIdentity (
        new [] {
          user.DisplayName == null ? null : new Claim (ClaimTypes.Name, user.DisplayName),
            user.PhotoUrl == null ? null : new Claim (ClaimTypes.Uri, user.PhotoUrl),
            new Claim (ClaimTypes.NameIdentifier, user.ID),
        }.Where (c => c != null).ToArray (),
        "Kumo User");
      claimsPrincipal = new ClaimsPrincipal (identity);
      idToken = user.IdToken;

      if (provider != null) {
        provider.NotifyAuthenticationStateChanged ();
      }
    }

    [JSInvokable]
    public static void MarkUserAsAnonymous () {
      claimsPrincipal = new ClaimsPrincipal (new ClaimsIdentity ());
      idToken = null;

      if (provider != null) {
        provider.NotifyAuthenticationStateChanged ();
      }
    }
  }
}
