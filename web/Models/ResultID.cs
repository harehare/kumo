using System;
using System.Text;

namespace Kumo.Models
{
    public class ResultID
    {
        private string id;

        private ResultID() { }

        public ResultID(string id)
        {
            this.id = id;
        }
        public ResultID(string url, string selector)
        {
            this.id = Convert.ToBase64String(Encoding.UTF8.GetBytes($"{url},{selector}"));
        }

        public override string ToString()
        {
            return id;
        }

        public (string Url, string Selector) Extract()
        {
            string[] tokens = Encoding.UTF8.GetString(Convert.FromBase64String(this.id)).Split(',');
            return (tokens[0], tokens[1]);
        }
    }
}
