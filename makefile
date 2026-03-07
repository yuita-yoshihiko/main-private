add-wt:
	wtp add -b ${WT_NAME}

rm-wt:
	wtp remove ${WT_NAME} --with-branch
