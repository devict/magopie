package com.bullercodeworks.magopie.adapters;

import android.content.Context;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.PopupMenu;
import android.widget.TextView;
import android.widget.Toast;

import com.bullercodeworks.magopie.R;
import com.bullercodeworks.magopie.State;

import org.w3c.dom.Text;

import java.util.List;

import go.magopie.Magopie;

public class ResultAdapter extends ArrayAdapter<Magopie.Torrent> {
  private Context ctx;
  private State state;
  public ResultAdapter(State _state, Context _ctx) {
    super(_ctx, R.layout.torrent_list_row, R.id.torrentTitle, _state.results);
    ctx = _ctx;
    state = _state;
  }

  public class ViewHolder {
    public TextView torrentTitle;
    public TextView seeds;
    public TextView leechers;
    public TextView site;

    public Magopie.Torrent torrent;

    ViewHolder(View row) {
      this.torrentTitle = (TextView)row.findViewById(R.id.torrentTitle);
      this.seeds = (TextView)row.findViewById(R.id.txtSeeds);
      this.leechers = (TextView)row.findViewById(R.id.txtLeechers);
      this.site = (TextView)row.findViewById(R.id.torrentSite);
    }
  }

  @Override
  public View getView(int position, View convertView, final ViewGroup parent) {
    View row = super.getView(position, convertView, parent);
    ViewHolder holder = (ViewHolder)row.getTag();
    if(holder == null) {
      holder = new ViewHolder(row);
      row.setTag(holder);
    }
    final Magopie.Torrent wrk = getItem(position);
    if(wrk != null) {
      holder.torrent = wrk;
      //holder.torrentFile.setText(wrk.getMagnetURI());
      holder.torrentTitle.setText(wrk.getTitle());
      holder.seeds.setText(String.valueOf(wrk.getSeeders()));
      holder.leechers.setText(String.valueOf(wrk.getLeechers()));
      holder.site.setText(state.sites.get(wrk.getSiteID()));
    }
    row.setOnLongClickListener(new View.OnLongClickListener() {
      @Override
      public boolean onLongClick(View view) {
        PopupMenu menu = new PopupMenu(getContext())
        return false;
      }
    });
    row.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        if(state.client.Download(((ViewHolder) view.getTag()).torrent)) {
          Toast.makeText(ctx, "Torrent download triggered", Toast.LENGTH_SHORT).show();
        } else {
          Toast.makeText(ctx, "Failed to trigger download!", Toast.LENGTH_SHORT).show();
        }
      }
    });
    return row;
  }
}