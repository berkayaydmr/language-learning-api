SELECT
	EXISTS (
		SELECT
			1
		FROM
			words
		WHERE
			id = ?
	);