namespace Kumo.Models
{
    public interface NotifyDestination
    {
        string Dest();
    }

    public class WebHookDestination
    {
        private string url;
        private WebHookDestination() { }
        public WebHookDestination(string url)
        {
            // TODO: check
            this.url = url;
        }

        public string Dest()
        {
            return this.url;
        }
    }

    public class SlackDestination
    {
        private string url;
        private SlackDestination() { }
        public SlackDestination(string url)
        {
            // TODO: check
            this.url = url;
        }

        public string Dest()
        {
            return this.url;
        }
    }

    public class EmailDestination
    {
        private string email;
        private EmailDestination() { }
        public EmailDestination(string email)
        {
            // TODO: check
            this.email = email;
        }

        public string Dest()
        {
            return this.email;
        }
    }
}
