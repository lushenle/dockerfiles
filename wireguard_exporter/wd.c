#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define WAIT_SECOND 2

long int getCurrentDownloadRates(long int *save_rate);

int main(int argc, char *argv[])
{
    long int start_download_rates;
    long int end_download_rates;
    while (1)
    {
        getCurrentDownloadRates(&start_download_rates);
        sleep(WAIT_SECOND);
        getCurrentDownloadRates(&end_download_rates);
        printf("download is : %.2lf Bytes/s\n", (float)(end_download_rates - start_download_rates) / WAIT_SECOND);
    }
    exit(EXIT_SUCCESS);
}

long int getCurrentDownloadRates(long int *save_rate)
{
    FILE *net_dev_file;
    char buffer[1024];
    size_t bytes_read;
    char *match;
    if ((net_dev_file = fopen("/proc/net/dev", "r")) == NULL)
    {
        printf("open file /proc/net/dev/ error!\n");
        exit(EXIT_FAILURE);
    }
    bytes_read = fread(buffer, 1, sizeof(buffer), net_dev_file);
    fclose(net_dev_file);
    if (bytes_read == 0)
    {
        exit(EXIT_FAILURE);
    }
    buffer[bytes_read] = '\0';
    match = strstr(buffer, "eth0:");
    if (match == NULL)
    {
        printf("no eth0 keyword to find!\n");
        exit(EXIT_FAILURE);
    }
    sscanf(match, "eth0:%ld", save_rate);
    return *save_rate;
}
