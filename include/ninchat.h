#ifndef NINCHAT_H
#define NINCHAT_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// CFFI BEGIN

typedef uintptr_t ninchat_session;

typedef struct {
	const void *data;
	size_t size;
	size_t _internal_private_reserved;
} ninchat_frame;

typedef void (*ninchat_session_event_callback)(
		void *context,
		const char *params,
		size_t params_len);

typedef void (*ninchat_event_callback)(
		void *context,
		const char *params,
		size_t params_len,
		const ninchat_frame payload[],
		unsigned int payload_len,
		bool last_reply);

typedef void (*ninchat_close_callback)(
		void *context);

typedef void (*ninchat_conn_state_callback)(
		void *context,
		const char *state);

typedef void (*ninchat_conn_active_callback)(
		void *context);

typedef void (*ninchat_log_callback)(
		void *context,
		const char *message,
		size_t message_len);

#ifndef NINCHAT_NO_PROTOTYPES

ninchat_session ninchat_session_new(void);

void ninchat_session_delete(ninchat_session s);

void ninchat_session_on_session_event(
		ninchat_session s,
		ninchat_session_event_callback cb,
		void *context);

void ninchat_session_on_event(
		ninchat_session s,
		ninchat_event_callback cb,
		void *context);

void ninchat_session_on_close(
		ninchat_session s,
		ninchat_close_callback cb,
		void *context);

void ninchat_session_on_conn_state(
		ninchat_session s,
		ninchat_conn_state_callback cb,
		void *context);

void ninchat_session_on_conn_active(
		ninchat_session s,
		ninchat_conn_active_callback cb,
		void *context);

void ninchat_session_on_log(
		ninchat_session s,
		ninchat_log_callback cb,
		void *context);

void ninchat_session_set_params(
		ninchat_session s,
		const char *params,
		size_t params_len);

void ninchat_session_open(ninchat_session s);

void ninchat_session_close(ninchat_session s);

char *ninchat_session_send(
		ninchat_session s,
		const char *params,
		size_t params_len,
		const ninchat_frame payload[],
		unsigned int payload_len,
		int64_t *action_id_ptr);

#endif

// CFFI END

#ifdef __cplusplus
}
#endif

#endif
