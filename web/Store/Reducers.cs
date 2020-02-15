using System.Linq;
using Blazor.Fluxor;
using Kumo.Models;

namespace Kumo.Store {
  public class FetchItemsCompleteReducer : Reducer<AppState, FetchItemsCompleteAction> {
    public override AppState Reduce (AppState state, FetchItemsCompleteAction action) {
      return new AppState (
        isProgress: false,
        errorMessage: state.ErrorMessage,
        items: action.Items);
    }
  }

  public class AddItemCompleteReducer : Reducer<AppState, AddItemCompleteAction> {
    public override AppState Reduce (AppState state, AddItemCompleteAction action) {
      return new AppState (
        isProgress: false,
        errorMessage: state.ErrorMessage,
        items: new Item[] { action.Item }.Cast<Item> ().Concat (state.Items).ToArray ());
    }
  }

  public class UpdateItemCompleteReducer : Reducer<AppState, UpdateItemCompleteAction> {
    public override AppState Reduce (AppState state, UpdateItemCompleteAction action) {
      return new AppState (
        isProgress: false,
        errorMessage: state.ErrorMessage,
        items: state.Items.Cast<Item> ().Select (item => item.Equals (action.Item) ? action.Item : item).ToArray ());
    }
  }

  public class DeleteItemCompleteReducer : Reducer<AppState, DeleteItemCompleteAction> {
    public override AppState Reduce (AppState state, DeleteItemCompleteAction action) {
      return new AppState (
        isProgress: false,
        errorMessage: state.ErrorMessage,
        items: state.Items.Cast<Item> ().Where (item => !item.Equals (action.Item)).ToArray ());
    }
  }
}
