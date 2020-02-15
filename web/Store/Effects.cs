using System;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading.Tasks;
using Blazor.Fluxor;
using Kumo.Models;
using Microsoft.AspNetCore.Components;
using Microsoft.JSInterop;

namespace Kumo.Store {
  public class FetchItemsEffect : Effect<FetchItemsAction> {
    private readonly HttpClient httpClient;

    public FetchItemsEffect (HttpClient httpClient) {
      this.httpClient = httpClient;
    }

    protected async override Task HandleAsync (FetchItemsAction action, IDispatcher dispatcher) {
      // HttpClient.DefaultRequestHeaders.Authorization =
      //   new AuthenticationHeaderValue("Bearer", action.User.IdToken);
      Item[] items = Array.Empty<Item> ();
      try {
        // TODO:
        items = await httpClient.GetJsonAsync<Item[]> ("http://localhost:8081/api/items");
        Console.WriteLine (items);
      } catch (Exception e) {
        // TODO: Should really dispatch an error action
        Console.WriteLine (e);
      }

      dispatcher.Dispatch (new FetchItemsCompleteAction (items));
    }
  }

  public class AddItemEffect : Effect<AddItemAction> {
    private readonly HttpClient httpClient;

    public AddItemEffect (HttpClient httpClient) {
      this.httpClient = httpClient;
    }

    protected async override Task HandleAsync (AddItemAction action, IDispatcher dispatcher) {
      try {
        var url = action.Item.Url;
        var selector = action.Item.Selector;
        await Task.WhenAll (new [] {
          httpClient.GetJsonAsync<Item> ($"http://localhost:8081/spider/entry?url={url}&selector={selector}"),
            httpClient.PostJsonAsync<Item> ("http://localhost:8081/api/items", action.Item)
        });
      } catch {
        // TODO: Should really dispatch an error action
      }
      dispatcher.Dispatch (new AddItemCompleteAction (action.Item));
    }
  }

  public class UpdateItemEffect : Effect<UpdateItemAction> {
    private readonly HttpClient httpClient;

    public UpdateItemEffect (HttpClient httpClient) {
      this.httpClient = httpClient;
    }

    protected async override Task HandleAsync (UpdateItemAction action, IDispatcher dispatcher) {
      try {
        var url = action.Item.Url;
        var selector = action.Item.Selector;
        await Task.WhenAll (new [] {
          httpClient.GetJsonAsync<Item> ($"http://localhost:8081/spider/entry?url={url}&selector={selector}"),
            httpClient.PostJsonAsync<Item> ("http://localhost:8081/api/items", action.Item)
        });
      } catch {
        // TODO: Should really dispatch an error action
      }
      dispatcher.Dispatch (new UpdateItemCompleteAction (action.Item));
    }
  }

  public class DeleteItemEffect : Effect<DeleteItemAction> {
    private readonly HttpClient httpClient;

    public DeleteItemEffect (HttpClient httpClient) {
      this.httpClient = httpClient;
    }

    protected async override Task HandleAsync (DeleteItemAction action, IDispatcher dispatcher) {
      // HttpClient.DefaultRequestHeaders.Authorization =
      //   new AuthenticationHeaderValue("Bearer", action.User.IdToken);
      try {
        // TODO:
        await httpClient.DeleteAsync ("test-data/items.json");
      } catch {
        // TODO: Should really dispatch an error action
      }
      dispatcher.Dispatch (new DeleteItemCompleteAction (action.Item));
    }
  }
}
