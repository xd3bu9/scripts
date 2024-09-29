import requests
import os
from argparse import ArgumentParser

# Define the base URL and dataset ID
base_url = "https://api.trickest.io/solutions/v1/public/solution/a7cba1f1-df07-4a5c-876a-953f178996be/view"
dataset_id = "937bfab9-7704-4954-b942-831bafbaf7d"
# key = print(os.environ['TRICKEST'])
def get_args():
    ArgumentParser.add_argument("-q", metavar="QUERY", help="query for the request")
    return ArgumentParser.parse_args()

# def fetch_all_results(base_url, dataset_id, limit=20, api_key=key):
#     all_results = []
#     offset = 0
#     headers = {"Authorization": f"Token ${api_key}"}

#     while True:
#         params = {"offset": offset, "limit": limit, "dataset_id": dataset_id, "q": f"{query}"}
#         response = requests.get(base_url, headers=headers, params=params)
#         # Raise an error for bad status codes
#         response.raise_for_status()

#         results = response.json().get("results", [])
#         if not results:
#             break

#         all_results.extend(results)
#         offset += limit

#     return all_results


if __name__ == "__main__":
    # if not os.environ.get('TRICKEST'):
    #     print(
    #         "Add TRICKEST to the system environment variables with the value of your api key."
    #     )
    #     exit(1)
    # else:
        # results = fetch_all_results(base_url, dataset_id, key)
        # print(results)
    print(os.getenv("TRICKEST"))
        