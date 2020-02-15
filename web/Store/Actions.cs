using Kumo.Models;

namespace Kumo.Store {
  public class FetchItemsAction {
    public readonly User User;

    public FetchItemsAction (User user) {
      User = user;
    }
  }

  public class FetchItemsCompleteAction {
    public readonly Item[] Items;

    public FetchItemsCompleteAction (Item[] items) {
      Items = items;
    }
  }

  public class AddItemAction {
    public readonly User User;
    public readonly Item Item;

    public AddItemAction (User user, Item item) {
      User = user;
      Item = item;
    }
  }

  public class AddItemCompleteAction {
    public readonly Item Item;

    public AddItemCompleteAction (Item item) {
      Item = item;
    }
  }

  public class UpdateItemAction {
    public readonly User User;
    public readonly Item Item;

    public UpdateItemAction (User user, Item item) {
      User = user;
      Item = item;
    }
  }

  public class UpdateItemCompleteAction {
    public readonly Item Item;

    public UpdateItemCompleteAction (Item item) {
      Item = item;
    }
  }

  public class DeleteItemAction {
    public readonly User User;
    public readonly Item Item;

    public DeleteItemAction (User user, Item item) {
      User = user;
      Item = item;
    }
  }

  public class DeleteItemCompleteAction {
    public readonly Item Item;

    public DeleteItemCompleteAction (Item item) {
      Item = item;
    }
  }
}
