using Blazor.Fluxor;

namespace Kumo.Store
{
    public class AppFeature : Feature<AppState>
    {
        public override string GetName() => "Kumo";
        protected override AppState GetInitialState() => new AppState();
    }
}
