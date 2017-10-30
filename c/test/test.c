#define _GNU_SOURCE

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <pthread.h>

#include "../../include/ninchat.h"

struct context {
	bool opened;
	int received_events;
	bool closed;
};

static pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t cond = PTHREAD_COND_INITIALIZER;

static void handle_session_event(
		void *context,
		const char *params,
		size_t paramslen)
{
	struct context *ctx = context;

	char buf[paramslen + 1];
	memset(buf, 0, paramslen + 1);
	memcpy(buf, params, paramslen);

	printf("handle_session_event: params: %s\n", buf);

	pthread_mutex_lock(&mutex);
	ctx->opened = true;
	pthread_cond_broadcast(&cond);
	pthread_mutex_unlock(&mutex);
}

static void handle_event(
		void *context,
		const char *params,
		size_t paramslen,
		const ninchat_frame payload[],
		unsigned int payloadlen,
		bool lastreply)
{
	struct context *ctx = context;

	char buf[paramslen + 1];
	memset(buf, 0, paramslen + 1);
	memcpy(buf, params, paramslen);

	printf("handle_event: params: %s\n", buf);
	printf("handle_event: payload length: %u\n", payloadlen);

	for (unsigned i = 0; i < payloadlen; i++) {
		printf("handle_event: payload part #%u size: %lu\n", i, payload[i].size);

		char buf[payload[i].size + 1];
		memset(buf, 0, payload[i].size + 1);
		memcpy(buf, payload[i].data, payload[i].size);

		printf("handle_event: payload part #%u data: %s\n", i, buf);
	}

	printf("handle_event: last reply: %s\n", lastreply ? "true" : "false");

	pthread_mutex_lock(&mutex);
	ctx->received_events++;
	pthread_cond_broadcast(&cond);
	pthread_mutex_unlock(&mutex);
}

static void handle_close(void *context)
{
	struct context *ctx = context;

	printf("handle_close\n");

	pthread_mutex_lock(&mutex);
	ctx->closed = true;
	pthread_cond_broadcast(&cond);
	pthread_mutex_unlock(&mutex);
}

static void handle_conn_state(void *context, const char *state)
{
	printf("handle_conn_state: state: %s\n", state);
}

static void handle_conn_active(void *context)
{
	printf("handle_conn_active\n");
}

static void handle_log(void *context, const char *message, size_t message_len)
{
	printf("handle_log: %s\n", message);
}

int main(int argc, char **argv)
{
	const char *session_params = "{ \"message_types\": [\"*\"] }";
	struct context ctx = {0};

	ninchat_session s = ninchat_session_new();
	ninchat_session_on_session_event(s, handle_session_event, &ctx);
	ninchat_session_on_event(s, handle_event, &ctx);
	ninchat_session_on_close(s, handle_close, &ctx);
	ninchat_session_on_conn_state(s, handle_conn_state, &ctx);
	ninchat_session_on_conn_active(s, handle_conn_active, &ctx);
	ninchat_session_on_log(s, handle_log, &ctx);
	ninchat_session_set_params(s, session_params, strlen(session_params));
	ninchat_session_open(s);

	pthread_mutex_lock(&mutex);
	while (!ctx.opened)
		pthread_cond_wait(&cond, &mutex);
	pthread_mutex_unlock(&mutex);

	int sent_events = 0;
	int64_t action_id;
	char *error;

	const char *params = "{ \"action\": \"describe_conn\" }";
	error = ninchat_session_send(s, params, strlen(params), NULL, 0, &action_id);
	if (error) {
		fprintf(stderr, "ninchat_session_send: error: %s\n", error);
		free(error);
	}

	printf("ninchat_session_send: action id: %ld\n", action_id);
	sent_events++;

	for (int i = 1; i < argc; i++) {
		const char *text = "{ \"text\": \"hello world\" }";
		ninchat_frame payload[] = {
			{text, strlen(text)},
			{"extra part", 10},
		};

		char *params;
		asprintf(&params, "{ \"action\": \"send_message\", \"user_id\": \"%s\", \"message_type\": \"ninchat.com/text\" }", argv[i]);
		error = ninchat_session_send(s, params, strlen(params), payload, 2, &action_id);
		free(params);
		if (error) {
			fprintf(stderr, "ninchat_session_send: error: %s\n", error);
			free(error);
		}

		printf("ninchat_session_send: action id: %ld\n", action_id);
		sent_events++;
	}

	pthread_mutex_lock(&mutex);
	while (ctx.received_events < sent_events)
		pthread_cond_wait(&cond, &mutex);
	pthread_mutex_unlock(&mutex);

	ninchat_session_close(s);

	pthread_mutex_lock(&mutex);
	while (!ctx.closed)
		pthread_cond_wait(&cond, &mutex);
	pthread_mutex_unlock(&mutex);

	ninchat_session_delete(s);

	printf("ok\n");
	return 0;
}
