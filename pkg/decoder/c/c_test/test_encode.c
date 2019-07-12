#include <stdio.h>
#include <sys/types.h>
#include <Rectangle.h>

static int write_out(const void *buffer, size_t size, void *app_key){
    FILE *out_fp = app_key;
    size_t wrote = fwrite(buffer, 1, size, out_fp);
    return (wrote == size) ? 0 : -1;
}

int main(int ac, char **av){
    Rectangle_t * rectangle;
    asn_enc_rval_t ec;

    rectangle = calloc(1, sizeof(Rectangle_t));

    if (!rectangle){
        perror("calloc() failed");
        exit(1);
    }

    rectangle->height = 42;
    rectangle-> width = 23;

    if (ac < 2){
        fprintf(stderr, "Specify filename for UPER output\n");
    } else {
        const char *filename = av[1];
        FILE *fp = fopen(filename, "wb");

        if (!fp) {
            perror(filename);
            exit(1);
        }

        ec = uper_encode(&asn_DEF_Rectangle, 0, rectangle, write_out, fp);
        fclose(fp);
        if (ec.encoded == -1){
            fprintf(stderr, "Could not encode Rectangle (at %s)\n",
                ec.failed_type ? ec.failed_type->name : "unknown");
            exit(1);
        } else {
            fprintf(stderr, "Created %s with UPER encoded Rectangle\n", filename);
        }
    }
    
    xer_fprint(stdout, &asn_DEF_Rectangle, rectangle);

    return 0;
}