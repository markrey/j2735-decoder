#include <stdio.h>
#include <sys/types.h>
#include <MessageFrame.h>

int main(int ac, char **av){
    char buf[1024];
    asn_dec_rval_t rval;
    MessageFrame_t *msgFrame = 0;
    FILE *fp;
    size_t size;
    char *filename;

    if(ac != 2) {
    fprintf(stderr, "Usage: %s <file.ber>\n", av[0]);
        exit(1);
    } else {
        filename = av[1];
    }

    /* Open input file as read-only binary */ 
    fp = fopen(filename, "rb");
    if(!fp) {
        perror(filename);
        exit(1); 
    }
    /* Read up to the buffer size */
    size = fread(buf, 1, sizeof(buf), fp);
    printf("bytes read: %ld\n", size);
    fclose(fp);
    if(!size) {
        fprintf(stderr, "%s: Empty or broken\n", filename);
        exit(1); 
    }
    /* Decode the input buffer as Rectangle type */
    rval = uper_decode_complete(0, &asn_DEF_MessageFrame, (void **)&msgFrame, buf, size);
    if(rval.code != RC_OK) {
        fprintf(stderr, "%s: Broken Rectangle encoding at byte %ld\n", filename,
            (long)rval.consumed);
        exit(1); 
    }
    printf("bytes consumed: %ld\n", rval.consumed);
    /* Print the decoded Rectangle type as XML */
    xer_fprint(stdout, &asn_DEF_MessageFrame, msgFrame);
    asn
    return 0; /* Decoding finished successfully */
}