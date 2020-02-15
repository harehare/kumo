using System;
using Kumo.Models;

namespace Kumo.Store {
  public class AppState {
    public bool IsProgress { get; private set; }
    public string ErrorMessage { get; private set; }
    public Item[] Items { get; private set; }

    public AppState () {
      IsProgress = false;
      ErrorMessage = null;
      Items = Array.Empty<Item> ();
    }

    public AppState (bool isProgress, string errorMessage, Item[] items) {
      IsProgress = IsProgress;
      ErrorMessage = errorMessage;
      Items = items;
    }
  }
}
