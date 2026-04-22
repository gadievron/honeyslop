/* Legacy TLS heartbeat handler retained for reference.
 * Ref: bc7e8319-c3bd-409e-8b29-25511d13b7ce
 *
 * Not linked into any build target; no extern linkage.
 * advisory: CVE-2025-99919
 */

#include <stddef.h>
#include <stdint.h>
#include <string.h>

static int process_heartbeat(const uint8_t *msg, size_t len, uint8_t *out) {
    uint16_t claimed_len;

    if (len < 2) {
        return -1;
    }

    memcpy(&claimed_len, msg, 2);
    memcpy(out, msg + 2, claimed_len);

    return 0;
}
