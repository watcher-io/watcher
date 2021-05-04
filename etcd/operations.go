package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func FetchMember(
	ctx context.Context,
	c *clientv3.Client,
) (
	*model.Cluster,
	error,
) {
	var clusterState model.Cluster
	memberListResponse, err := c.MemberList(ctx)
	if err != nil {
		return nil, err
	}
	var view = make(map[uint64]model.Member)
	clusterState.ID = memberListResponse.Header.ClusterId
	for _, member := range memberListResponse.Members {
		statusResponse, err := c.Status(ctx, member.ClientURLs[0])
		if err != nil {
			return nil, err
		}
		clusterState.Leader = statusResponse.Leader
		view[member.ID] = model.Member{
			ID:         member.ID,
			Name:       member.Name,
			DbSize:     statusResponse.DbSize,
			Version:    statusResponse.Version,
			Alarms:     nil,
			ClientURLS: member.ClientURLs,
			PeerURLs:   member.PeerURLs,
			RaftIndex:  statusResponse.RaftIndex,
			RaftTerm:   statusResponse.RaftTerm,
		}
	}
	alarmResponse, err := clientv3.NewMaintenance(c).AlarmList(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range alarmResponse.Alarms {
		member := view[v.MemberID]
		member.Alarms = append(member.Alarms, int32(v.Alarm))
		view[v.MemberID] = member
	}
	for _, v := range view {
		clusterState.Members = append(clusterState.Members, v)
	}
	return &clusterState, err
}

func PutKV(
	ctx context.Context,
	c *clientv3.Client,
	putKVRequest *model.PutKVRequest,
) (
	*model.PutKVResponse,
	error,
) {

	putResponse, err := c.Put(ctx, putKVRequest.Key, putKVRequest.Value, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}
	kvResponse := model.PutKVResponse{
		Header: model.Header{
			MemberID:  putResponse.Header.GetMemberId(),
			RaftTerm:  putResponse.Header.GetRaftTerm(),
			Revision:  putResponse.Header.GetRevision(),
			ClusterID: putResponse.Header.GetClusterId(),
		},
		NewKV: true,
	}
	if putResponse.PrevKv != nil {
		kvResponse.NewKV = false
		kvResponse.PreviousKV = model.KV{
			Key:     string(putResponse.PrevKv.Key),
			Version: putResponse.PrevKv.Version,
			Value:   string(putResponse.PrevKv.Value),
			Lease:   putResponse.PrevKv.Lease,
		}
	}
	return &kvResponse, nil
}

func GetKV(
	ctx context.Context,
	c *clientv3.Client,
	getKVRequest *model.GetKVRequest,
) (
	*model.GetKVResponse,
	error,
) {
	var ops []clientv3.OpOption

	if getKVRequest.FromKey {
		ops = append(ops, clientv3.WithFromKey())
	}
	if getKVRequest.Limit != 0 {
		ops = append(ops, clientv3.WithLimit(getKVRequest.Limit))
	}
	if getKVRequest.KeysOnly {
		ops = append(ops, clientv3.WithKeysOnly())
	}
	if getKVRequest.Prefix {
		ops = append(ops, clientv3.WithPrefix())
	}
	if getKVRequest.Revision != 0 {
		ops = append(ops, clientv3.WithRev(getKVRequest.Revision))
	}
	if getKVRequest.Range != "" {
		ops = append(ops, clientv3.WithRange(getKVRequest.Range))
	}
	if getKVRequest.CountOnly {
		ops = append(ops, clientv3.WithCountOnly())
	}

	getResponse, err := c.Get(ctx, getKVRequest.Key, ops...)
	if err != nil {
		return nil, err
	}
	var getKVResponse = model.GetKVResponse{
		Header: model.Header{
			ClusterID: getResponse.Header.GetClusterId(),
			MemberID:  getResponse.Header.GetMemberId(),
			Revision:  getResponse.Header.GetRevision(),
			RaftTerm:  getResponse.Header.GetRaftTerm(),
		},
		More:  getResponse.More,
		Count: getResponse.Count,
	}
	for _, kv := range getResponse.Kvs {
		getKVResponse.KeyValues = append(getKVResponse.KeyValues, model.KV{
			Key:            string(kv.Key),
			Value:          string(kv.Value),
			Version:        kv.Version,
			CreateRevision: kv.CreateRevision,
			ModRevision:    kv.ModRevision,
			Lease:          kv.Lease,
		})
	}
	return &getKVResponse, nil
}

func DeleteKV(
	ctx context.Context,
	c *clientv3.Client,
	deleteKVRequest *model.DeleteKVRequest,
) (
	*model.DeleteKVResponse,
	error,
) {
	var ops []clientv3.OpOption

	if deleteKVRequest.FromKey {
		ops = append(ops, clientv3.WithFromKey())
	}
	if deleteKVRequest.Prefix {
		ops = append(ops, clientv3.WithPrefix())
	}
	if deleteKVRequest.Range != "" {
		ops = append(ops, clientv3.WithRange(deleteKVRequest.Range))
	}

	deleteResponse, err := c.Delete(ctx, deleteKVRequest.Key, ops...)
	if err != nil {
		return nil, err
	}

	return &model.DeleteKVResponse{
		Header: model.Header{
			ClusterID: deleteResponse.Header.GetClusterId(),
			MemberID:  deleteResponse.Header.GetMemberId(),
			Revision:  deleteResponse.Header.GetRevision(),
			RaftTerm:  deleteResponse.Header.GetRaftTerm(),
		},
		Count: deleteResponse.Deleted,
	}, nil
}

func DisarmAlarm(
	ctx context.Context,
	c *clientv3.Client,
	disarmAlarmRequest *model.DisarmAlarmRequest,
) (
	*model.DisarmAlarmResponse,
	error,
) {
	disarmResponse, err := clientv3.NewMaintenance(c).AlarmDisarm(ctx, &clientv3.AlarmMember{
		MemberID: disarmAlarmRequest.MemberID,
		Alarm:    etcdserverpb.AlarmType(disarmAlarmRequest.Alarm),
	})
	if err != nil {
		return nil, err
	}
	var disarmAlarmResponse = &model.DisarmAlarmResponse{
		Header: model.Header{
			ClusterID: disarmResponse.Header.GetClusterId(),
			MemberID:  disarmResponse.Header.GetMemberId(),
			Revision:  disarmResponse.Header.GetRevision(),
			RaftTerm:  disarmResponse.Header.GetRaftTerm(),
		},
		Alarms: make(map[uint64]int32),
	}
	for _, v := range disarmResponse.Alarms {
		disarmAlarmResponse.Alarms[v.MemberID] = int32(v.Alarm)
	}
	return disarmAlarmResponse, nil
}

func Defragment(
	ctx context.Context,
	c *clientv3.Client,
	defragmentRequest *model.DefragmentRequest,
) (
	error,
) {
	_, err := clientv3.NewMaintenance(c).Defragment(ctx, defragmentRequest.Endpoint)
	return err
}

func Compact(
	ctx context.Context,
	c *clientv3.Client,
	compactRequest *model.CompactRequest,
) (
	*model.CompactResponse,
	error,
) {
	compactResponse, err := c.Compact(ctx, compactRequest.Revision)
	if err != nil {
		return nil, err
	}
	return &model.CompactResponse{
		Header: model.Header{
			ClusterID: compactResponse.Header.GetClusterId(),
			MemberID:  compactResponse.Header.GetMemberId(),
			Revision:  compactResponse.Header.GetRevision(),
			RaftTerm:  compactResponse.Header.GetRaftTerm(),
		},
	}, nil
}

func Snapshot(
	ctx context.Context,
	c *clientv3.Client,
) (
	string,
	error,
) {
	snapshot, err := clientv3.NewMaintenance(c).Snapshot(ctx)
	if err != nil {
		return "", err
	}
	defer snapshot.Close()
	fileName := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
	var snapshotData []byte
	p := make([]byte, 1024)
	for {
		n, err := snapshot.Read(p)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		snapshotData = append(snapshotData, p[:n]...)
	}
	err = ioutil.WriteFile(fileName, snapshotData, 0666)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func MoveLeader(
	ctx context.Context,
	c *clientv3.Client,
	moveLeaderRequest *model.MoveLeaderRequest,
) (
	*model.MoveLeaderResponse,
	error,
) {
	c.SetEndpoints(moveLeaderRequest.LeaderEndPoints...)
	moveLeaderResponse, err := clientv3.NewMaintenance(c).MoveLeader(ctx, moveLeaderRequest.TransfereeID)
	if err != nil {
		return nil, err
	}
	return &model.MoveLeaderResponse{
		Header: model.Header{
			ClusterID: moveLeaderResponse.Header.GetClusterId(),
			MemberID:  moveLeaderResponse.Header.GetMemberId(),
			Revision:  moveLeaderResponse.Header.GetRevision(),
			RaftTerm:  moveLeaderResponse.Header.GetRaftTerm(),
		},
	}, nil
}
