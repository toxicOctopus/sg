package database

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/toxicOctopus/sg/internal/twitch"
)

func GetRegisteredChannels(ctx context.Context, db *pgx.Conn) (twitch.RegisteredChannels, error) {
	rc := twitch.RegisteredChannels{}
	rows, err := db.Query(ctx, `
	select
		c.id,
		c.twitch_name,
		c.action_cooldown,
		array_to_json(ARRAY(select json_build_object(
			'id', ce.emote_id,
			'name', e.name,
			'image_path', e.image_path,
			'action_type', ce.action_type_id,
			'action_name', a.name,
			'action_source', a.action_source
		) :: jsonb
			from 
				channel_emotes ce
			left join 
				emotes e on ce.emote_id = e.id
			left join
				action_types a on ce.action_type_id = a.id
			where 
				ce.channel_id = c.id
		)) as emotes
	from
		channels c`)
	if err != nil {
		return rc, errors.Wrap(err, "getting channels from db")
	}

	for rows.Next() {
		channel := twitch.Channel{}
		var emotesJson []byte
		emoteList := twitch.EmoteList{}
		var actionCD pgtype.Int8
		err = rows.Scan(&channel.ID, &channel.Name, &actionCD, &emotesJson)
		if err != nil {
			return rc, errors.Wrap(err, "scanning channel rows")
		}
		actionCDDuration, err := time.ParseDuration(strconv.Itoa(int(actionCD.Int)) + "s")
		if err == nil {
			channel.ActionCD = actionCDDuration
		}

		if len(emotesJson) > 0 {
			err = json.Unmarshal(emotesJson, &emoteList)
			if err != nil {
				return rc, errors.Wrap(err, "unmarshalling emote object")
			}
			channel.Emotes = emoteList
		}

		rc = append(rc, channel)
	}

	return rc, nil
}
