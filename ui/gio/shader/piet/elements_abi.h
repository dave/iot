// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.

struct elements_descriptor_set_layout {
	struct buffer_descriptor binding0;
	struct buffer_descriptor binding1;
	struct buffer_descriptor binding2;
	struct buffer_descriptor binding3;
};

extern coroutine elements_coroutine_begin(struct program_data *data,
	int32_t workgroupX, int32_t workgroupY, int32_t workgroupZ,
	void *workgroupMemory,
	int32_t firstSubgroup,
	int32_t subgroupCount) ATTR_HIDDEN;

extern bool elements_coroutine_await(coroutine r, yield_result *res) ATTR_HIDDEN;
extern void elements_coroutine_destroy(coroutine r) ATTR_HIDDEN;

extern const struct program_info elements_program_info ATTR_HIDDEN;
