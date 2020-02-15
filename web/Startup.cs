using Blazor.Fluxor;
using Kumo.Auth;
using Microsoft.AspNetCore.Components.Authorization;
using Microsoft.AspNetCore.Components.Builder;
using Microsoft.Extensions.DependencyInjection;

namespace Kumo {
  public class Startup {
    public void ConfigureServices (IServiceCollection services) {
      services.AddAuthorizationCore ();
      services.AddFluxor (options => {
        options.UseDependencyInjection (typeof (Startup).Assembly);
        options.AddMiddleware<Blazor.Fluxor.ReduxDevTools.ReduxDevToolsMiddleware> ();
      });
      services.AddScoped<AuthenticationStateProvider, ApiAuthenticationStateProvider> ();
    }

    public void Configure (IComponentsApplicationBuilder app) {
      app.AddComponent<App> ("app");
    }
  }
}
