package com.bullercodeworks.magopie.adapters;

import android.content.Context;
import android.content.Intent;
import android.util.Log;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.ImageView;
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
    // If state doesn't have sites at this point, try to grab them again.
    if(state.sites.isEmpty()) {
      state.UpdateSites();
    }
  }

  public class ViewHolder {
    public ImageView magnet;
    public TextView torrentTitle;
    public TextView seeds;
    public TextView leechers;
    public TextView site;

    public Magopie.Torrent torrent;

    ViewHolder(View row) {
      this.magnet = (ImageView)row.findViewById(R.id.magnetLink);
      this.torrentTitle = (TextView)row.findViewById(R.id.torrentTitle);
      this.seeds = (TextView)row.findViewById(R.id.txtSeeds);
      this.leechers = (TextView)row.findViewById(R.id.txtLeechers);
      this.site = (TextView)row.findViewById(R.id.torrentSite);
    }
  }

  @Override
  public View getView(int position, View convertView, final ViewGroup parent) {
    final View row = super.getView(position, convertView, parent);
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
        PopupMenu menu = new PopupMenu(ctx.getApplicationContext(), view);
        menu.getMenuInflater().inflate(R.menu.popup_menu, menu.getMenu());
        menu.setOnMenuItemClickListener(new PopupMenu.OnMenuItemClickListener() {
          public boolean onMenuItemClick(MenuItem item) {
            int id = item.getItemId();
            Magopie.Torrent torrent = ((ViewHolder)row.getTag()).torrent;
            if (id == R.id.download) {
              if (state.client.Download(torrent)) {
                Toast.makeText(ctx, "Torrent download triggered", Toast.LENGTH_SHORT).show();
              } else {
                Toast.makeText(ctx, "Failed to trigger download!", Toast.LENGTH_SHORT).show();
              }
              return true;
            } else if(id == R.id.send_to) {
              Intent sendIntent = new Intent();
              sendIntent.putExtra(Intent.EXTRA_SUBJECT, torrent.getTitle());
              sendIntent.putExtra(Intent.EXTRA_TEXT, wrk.getMagnetURI());
              sendIntent.setAction(Intent.ACTION_SEND);
              sendIntent.setType("text/plain");
              ctx.startActivity(Intent.createChooser(sendIntent, "Share Torrent!"));
              return true;
            }
            return false;
          }
        });
        menu.show();
        return true;
      }
    });
    holder.magnet.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        Magopie.Torrent torrent = ((ViewHolder)row.getTag()).torrent;
        Intent sendIntent = new Intent();
        sendIntent.putExtra(Intent.EXTRA_SUBJECT, torrent.getTitle());
        sendIntent.putExtra(Intent.EXTRA_TEXT, wrk.getMagnetURI());
        sendIntent.setAction(Intent.ACTION_SEND);
        sendIntent.setType("text/plain");
        ctx.startActivity(Intent.createChooser(sendIntent, "Share Torrent!"));
      }
    });
    return row;
  }
}